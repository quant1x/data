use std::any::Any;
use std::sync::mpsc::{sync_channel, Receiver, SyncSender};
use std::sync::{Arc, Mutex, Once};
use std::thread;
use std::time::Duration;

use anyhow::{anyhow, Context, Result};
use log::warn;

use crate::quotes::base_client::TcpClient;
use crate::quotes::base_consts::{QuotesError, CONN_TIMEOUT, POOL_MAX, RECV_TIMEOUT};
use crate::quotes::bestip::{get_fast_host, HOST_HQ};
use crate::quotes::cmd_company_info_category::{
    new_company_info_category, CompanyInfoCategory, CompanyInfoCategoryRequest,
};
use crate::quotes::cmd_company_info_content::{
    new_company_info_content, CompanyInfoContent, CompanyInfoContentRequest,
};
use crate::quotes::cmd_finance_info::{new_finance_info, FinanceInfo, FinanceInfoRequest};
use crate::quotes::cmd_heartbeat::{new_heartbeat, HeartBeatReply};
use crate::quotes::cmd_hello1::{new_hello1, Hello1Reply};
use crate::quotes::cmd_hello2::{new_hello2, Hello2Reply};
use crate::quotes::cmd_minute_time_data::{
    new_history_minute_time, new_minute_time, MinuteTimeReply, MinuteTimeRequest,
    HistoryMinuteTimeRequest,
};
use crate::quotes::cmd_transaction_data::{
    new_history_transaction, new_transaction, TransactionReply, TransactionRequest,
    HistoryTransactionRequest, TRANSACTION_MAX,
};
use crate::quotes::cmd_xdxr_info::{
    new_xdxr_info, XdxrInfo, XdxrInfoRequest,
};
use crate::quotes::cmd_security_bars::{
    new_index_bars, new_security_bars, SecurityBarsReply, SecurityBarsRequest, SECURITY_BARS_MAX,
};
use crate::quotes::cmd_security_count::{
    new_security_count, SecurityCountReply, SecurityCountRequest,
};
use crate::quotes::cmd_security_list::{new_security_list, SecurityListReply, SecurityListRequest};
use crate::quotes::cmd_security_quotes::{
    new_security_quotes, SecurityQuotesReply, SecurityQuotesRequest, SECURITY_QUOTES_MAX,
};
use crate::quotes::cmd_security_quotes_v2::{
    new_security_quotes_v2, V2SecurityQuotesReply, V2SecurityQuotesRequest, V2Stock,
    SECURITY_QUOTES_MAX_V2,
};
use crate::quotes::cmd_security_snapshot::{ExchangeState, Snapshot};
use crate::quotes::base_pool::ConnPool;
use crate::quotes::message::Message;
use crate::quotes::options::Options;
use crate::quotes::server::Server;
use crate::quotes::Stock;

pub struct StdApi {
    conn_pool: ConnPool,
    #[allow(dead_code)]
    opt: Options,
    queue: Arc<ServerQueue>,
}

struct ServerQueue {
    once: Once,
    servers: Vec<Server>,
    sender: Arc<SyncSender<Server>>,
    receiver: Mutex<Receiver<Server>>,
}

impl ServerQueue {
    fn new(servers: Vec<Server>, capacity: usize) -> Self {
        let cap = capacity.max(1);
        let (sender, receiver) = sync_channel(cap);
        Self {
            once: Once::new(),
            servers,
            sender: Arc::new(sender),
            receiver: Mutex::new(receiver),
        }
    }

    fn ensure_initialized(&self) {
        let sender = Arc::clone(&self.sender);
        self.once.call_once(|| {
            for server in &self.servers {
                if let Err(err) = sender.send(server.clone()) {
                    warn!("初始化服务器地址失败: {err}");
                }
            }
        });
    }

    fn acquire(&self) -> Option<Server> {
        self.ensure_initialized();
        let receiver = self.receiver.lock().expect("server receiver poisoned");
        receiver
            .recv()
            .ok()
            .filter(|srv| !srv.host.is_empty() && srv.port != 0)
    }

    fn release(&self, server: &Server) {
        if server.host.is_empty() || server.port == 0 {
            return;
        }
        self.ensure_initialized();
        if let Err(err) = self.sender.send(server.clone()) {
            warn!("返还服务器地址失败: {err}");
        }
    }

    fn len(&self) -> usize {
        self.servers.len()
    }

    fn release_callback(queue: &Arc<Self>) -> Arc<dyn Fn(&Server) + Send + Sync> {
        let weak = Arc::downgrade(queue);
        Arc::new(move |server: &Server| {
            if let Some(queue) = weak.upgrade() {
                queue.release(server);
            }
        })
    }
}

pub fn new_std_api() -> Result<StdApi> {
    let servers = get_fast_host(HOST_HQ)?;
    new_std_api_with_servers(servers)
}

pub fn new_std_api_with_servers(servers: Vec<Server>) -> Result<StdApi> {
    if servers.is_empty() {
        return Err(anyhow!("no available hosts"));
    }

    let server_count = servers.len();
    let mut max_cap = POOL_MAX.min(server_count);
    if max_cap == 0 {
        max_cap = 1;
    }

    let mut max_idle = max_cap;
    if let Ok(parallelism) = thread::available_parallelism() {
        let half = (parallelism.get() / 2).max(1);
        if max_idle > half {
            max_idle = half;
        }
    }

    let queue = Arc::new(ServerQueue::new(servers, POOL_MAX));

    let mut opt = Options::default();
    opt.connection_timeout = Duration::from_secs(CONN_TIMEOUT);
    opt.read_timeout = Duration::from_secs(RECV_TIMEOUT);
    opt.write_timeout = Duration::from_secs(RECV_TIMEOUT);
    opt.release_address = Some(ServerQueue::release_callback(&queue));

    let opt_shared = Arc::new(opt.clone());
    let queue_for_factory = Arc::clone(&queue);
    let factory_opt = Arc::clone(&opt_shared);
    let factory = move || -> Result<TcpClient> {
        let server = queue_for_factory
            .acquire()
            .ok_or_else(|| anyhow!(QuotesError::InvalidServerAddress))?;
        let mut client = TcpClient::new(factory_opt.as_ref().clone())?;
        if let Err(err) = client.connect(&server) {
            queue_for_factory.release(&server);
            return Err(err);
        }
        if let Err(err) = StdApi::tdx_hello1(&mut client) {
            let _ = client.close();
            return Err(err);
        }
        if let Err(err) = StdApi::tdx_hello2(&mut client) {
            let _ = client.close();
            return Err(err);
        }
        Ok(client)
    };

    let close_fn = move |client: TcpClient| -> Result<()> {
        client.close()
    };

    let ping_fn = move |client: &mut TcpClient| -> Result<()> {
        StdApi::tdx_ping(client)
    };

    let conn_pool = ConnPool::new(
        max_cap,
        max_idle,
        factory,
        close_fn,
        Some(ping_fn),
    )?;

    Ok(StdApi {
        conn_pool,
        opt,
        queue,
    })
}

impl StdApi {
    pub fn close(&self) {
        self.conn_pool.close_all();
    }

    /// TDX协议握手1 - 用于连接建立时的协议验证
    pub fn tdx_hello1(client: &mut TcpClient) -> Result<()> {
        let mut msg = new_hello1()?;
        client.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Hello1Reply>()
            .ok_or_else(|| anyhow!("invalid hello1 reply type"))?;
        log::info!("tdx-hello1: {}", reply.info);
        Ok(())
    }

    /// TDX协议握手2 - 用于连接建立时的协议验证
    pub fn tdx_hello2(client: &mut TcpClient) -> Result<()> {
        let mut msg = new_hello2()?;
        client.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Hello2Reply>()
            .ok_or_else(|| anyhow!("invalid hello2 reply type"))?;
        log::info!("tdx-hello2: {}", reply.info);
        Ok(())
    }

    /// TDX ping - 用于连接健康检查
    pub fn tdx_ping(client: &mut TcpClient) -> Result<()> {
        let mut msg = new_heartbeat()?;
        client.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<HeartBeatReply>()
            .ok_or_else(|| anyhow!("invalid heartbeat reply type"))?;
        if reply.is_empty() {
            return Err(anyhow!("ping failed: empty reply"));
        }
        Ok(())
    }

    pub fn command(&self, msg: &mut dyn Message) -> Result<()> {
        let conn = self.conn_pool.get_conn()?;
        let result = conn.command(msg);

        match result {
            Ok(()) => self.conn_pool.return_conn(conn),
            Err(err) => {
                let close_result = self.conn_pool.close_conn(conn);
                if let Err(close_err) = close_result {
                    warn!("close conn failed: {close_err}");
                }
                Err(err)
            }
        }
    }

    pub fn len(&self) -> usize {
        self.queue.len()
    }

    pub fn num_of_servers(&self) -> usize {
        self.queue.len()
    }

    pub fn get_max_idle_count(&self) -> usize {
        self.conn_pool.get_max_idle_count()
    }

    pub fn acquire_address(&self) -> Option<Server> {
        self.queue.acquire()
    }

    pub fn release_address(&self, server: &Server) {
        self.queue.release(server);
    }

    pub fn hello1(&self) -> Result<Hello1Reply> {
        let mut msg = new_hello1()?;
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Hello1Reply>()
            .ok_or_else(|| anyhow!("invalid hello1 reply type"))?;
        Ok(reply.clone())
    }

    pub async fn hello2(&self) -> Result<Hello2Reply> {
        let mut msg = new_hello2()?;
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Hello2Reply>()
            .ok_or_else(|| anyhow!("invalid hello2 reply type"))?;
        Ok(reply.clone())
    }

    pub fn heartbeat(&self) -> Result<HeartBeatReply> {
        let mut msg = new_heartbeat()?;
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<HeartBeatReply>()
            .ok_or_else(|| anyhow!("invalid heartbeat reply type"))?;
        Ok(reply.clone())
    }

    pub fn security_count(&self, market: u16) -> Result<SecurityCountReply> {
        let mut msg = new_security_count()?;
        msg.set_params(SecurityCountRequest { market });
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<SecurityCountReply>()
            .ok_or_else(|| anyhow!("invalid security count reply type"))?;
        Ok(reply.clone())
    }

    pub fn security_list(&self, market: u16, start: u16) -> Result<SecurityListReply> {
        let mut msg = new_security_list()?;
        msg.set_params(SecurityListRequest { market, start });
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<SecurityListReply>()
            .ok_or_else(|| anyhow!("invalid security list reply type"))?;
        Ok(reply.clone())
    }

    pub fn company_info_category(
        &self,
        market: u16,
        code: [u8; 6],
    ) -> Result<Vec<CompanyInfoCategory>> {
        let mut msg = new_company_info_category()?;
        msg.set_params(CompanyInfoCategoryRequest {
            market,
            code,
            unknown: 0,
        });
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Vec<CompanyInfoCategory>>()
            .ok_or_else(|| anyhow!("invalid company info category reply type"))?;
        Ok(reply.clone())
    }

    pub fn company_info_content(
        &self,
        request: CompanyInfoContentRequest,
    ) -> Result<CompanyInfoContent> {
        let mut msg = new_company_info_content()?;
        msg.set_params(request);
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<CompanyInfoContent>()
            .ok_or_else(|| anyhow!("invalid company info content reply type"))?;
        Ok(reply.clone())
    }

    pub fn security_quotes(&self, stocks: Vec<Stock>) -> Result<SecurityQuotesReply> {
        if stocks.is_empty() {
            return Err(anyhow!("security quotes requires at least one stock"));
        }
        if stocks.len() > SECURITY_QUOTES_MAX {
            return Err(anyhow!(
                "security quotes exceeds max {}",
                SECURITY_QUOTES_MAX
            ));
        }
        let mut msg = new_security_quotes()?;
        msg.set_params(SecurityQuotesRequest { stock_list: stocks });
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<SecurityQuotesReply>()
            .ok_or_else(|| anyhow!("invalid security quotes reply type"))?;
        Ok(reply.clone())
    }

    pub fn security_quotes_v2(&self, stocks: Vec<V2Stock>) -> Result<V2SecurityQuotesReply> {
        if stocks.is_empty() {
            return Err(anyhow!("security quotes v2 requires at least one stock"));
        }
        if stocks.len() > SECURITY_QUOTES_MAX_V2 {
            return Err(anyhow!(
                "security quotes v2 exceeds max {}",
                SECURITY_QUOTES_MAX_V2
            ));
        }
        let mut msg = new_security_quotes_v2()?;
        msg.set_params(V2SecurityQuotesRequest { stock_list: stocks });
        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<V2SecurityQuotesReply>()
            .ok_or_else(|| anyhow!("invalid security quotes v2 reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_snapshot(&self, codes: Vec<String>) -> Result<Vec<Snapshot>> {
        use crate::exchange::symbol::detect_market;
        use crate::exchange::symbol::get_security_code;

        if codes.is_empty() {
            return Err(anyhow!("get_snapshot requires at least one code"));
        }
        if codes.len() > SECURITY_QUOTES_MAX {
            return Err(anyhow!("get_snapshot exceeds max {}", SECURITY_QUOTES_MAX));
        }

        // Convert codes to stocks
        let mut stocks = Vec::new();
        for code in &codes {
            let (market_id, _, symbol) = detect_market(code);
            if symbol.len() == 6 {
                stocks.push(Stock {
                    market: market_id as u8,
                    code: symbol.clone(),
                });
            }
        }

        if stocks.is_empty() {
            return Err(anyhow!("no valid codes found"));
        }

        // Get security quotes
        let reply = self.security_quotes(stocks)?;

        // Convert to snapshots
        let mut snapshots = Vec::new();
        let current_date = "2025-10-19".to_string(); // TODO: Get current trading date

        for stock in &reply.list {
            let mut snapshot = Snapshot::default();
            // Copy basic fields from stock quote
            snapshot.security_code = get_security_code(stock.market as u16, &stock.code);
            snapshot.market = stock.market;
            snapshot.code = stock.code.clone();
            snapshot.price = stock.price;
            snapshot.last_close = stock.last_close;
            snapshot.open = stock.open;
            snapshot.high = stock.high;
            snapshot.low = stock.low;
            snapshot.vol = stock.vol;
            snapshot.amount = stock.amount;
            snapshot.date = current_date.clone();
            snapshot.active = stock.active1;
            snapshot.state = stock.state;

            // Set exchange state based on trading status
            // TODO: Implement real-time trading status check
            snapshot.exchange_state = ExchangeState::Normal; // Default to normal for now

            // Copy bid/ask data
            snapshot.bid1 = stock.bid1;
            snapshot.ask1 = stock.ask1;
            snapshot.bid_vol1 = stock.bid_vol1;
            snapshot.ask_vol1 = stock.ask_vol1;
            snapshot.bid2 = stock.bid2;
            snapshot.ask2 = stock.ask2;
            snapshot.bid_vol2 = stock.bid_vol2;
            snapshot.ask_vol2 = stock.ask_vol2;
            snapshot.bid3 = stock.bid3;
            snapshot.ask3 = stock.ask3;
            snapshot.bid_vol3 = stock.bid_vol3;
            snapshot.ask_vol3 = stock.ask_vol3;
            snapshot.bid4 = stock.bid4;
            snapshot.ask4 = stock.ask4;
            snapshot.bid_vol4 = stock.bid_vol4;
            snapshot.ask_vol4 = stock.ask_vol4;
            snapshot.bid5 = stock.bid5;
            snapshot.ask5 = stock.ask5;
            snapshot.bid_vol5 = stock.bid_vol5;
            snapshot.ask_vol5 = stock.ask_vol5;

            snapshots.push(snapshot);
        }

        Ok(snapshots)
    }

    pub fn get_finance_info(&self, code: &str) -> Result<FinanceInfo> {
        use crate::exchange::symbol::detect_market;

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_finance_info()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = FinanceInfoRequest {
            count: 1,
            market: market_id as u8,
            code: code_bytes,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<FinanceInfo>()
            .ok_or_else(|| anyhow!("invalid finance info reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_kline(&self, code: &str, category: u16, start: u16, count: u16) -> Result<SecurityBarsReply> {
        use crate::exchange::symbol::detect_market;

        if count as usize > SECURITY_BARS_MAX {
            return Err(anyhow!("kline count exceeds max {}", SECURITY_BARS_MAX));
        }

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_security_bars()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = SecurityBarsRequest {
            market: market_id as u16,
            code: code_bytes,
            category,
            i: 0,
            start,
            count,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<SecurityBarsReply>()
            .ok_or_else(|| anyhow!("invalid security bars reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_index_bars(&self, code: &str, category: u16, start: u16, count: u16) -> Result<SecurityBarsReply> {
        use crate::exchange::symbol::detect_market;

        if count as usize > SECURITY_BARS_MAX {
            return Err(anyhow!("index bars count exceeds max {}", SECURITY_BARS_MAX));
        }

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_index_bars()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = SecurityBarsRequest {
            market: market_id as u16,
            code: code_bytes,
            category,
            start,
            count,
            i: 0, // Will be set to 1 in set_params
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<SecurityBarsReply>()
            .ok_or_else(|| anyhow!("invalid index bars reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_minute_time(&self, code: &str, date: u32) -> Result<MinuteTimeReply> {
        use crate::exchange::symbol::detect_market;

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_minute_time()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = MinuteTimeRequest {
            market: market_id as u16,
            code: code_bytes,
            date,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<MinuteTimeReply>()
            .ok_or_else(|| anyhow!("invalid minute time reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_history_minute_time(&self, code: &str, date: u32) -> Result<MinuteTimeReply> {
        use crate::exchange::symbol::detect_market;

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_history_minute_time()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = HistoryMinuteTimeRequest {
            date,
            market: market_id as u8,
            code: code_bytes,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<MinuteTimeReply>()
            .ok_or_else(|| anyhow!("invalid history minute time reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_transaction(&self, code: &str, start: u16, count: u16) -> Result<TransactionReply> {
        use crate::exchange::symbol::detect_market;

        if count as usize > TRANSACTION_MAX {
            return Err(anyhow!("transaction count exceeds max {}", TRANSACTION_MAX));
        }

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_transaction()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = TransactionRequest {
            market: market_id as u16,
            code: code_bytes,
            start,
            count,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<TransactionReply>()
            .ok_or_else(|| anyhow!("invalid transaction reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_history_transaction(&self, code: &str, date: u32, start: u16, count: u16) -> Result<TransactionReply> {
        use crate::exchange::symbol::detect_market;

        if count as usize > TRANSACTION_MAX {
            return Err(anyhow!("history transaction count exceeds max {}", TRANSACTION_MAX));
        }

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_history_transaction()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = HistoryTransactionRequest {
            date,
            market: market_id as u16,
            code: code_bytes,
            start,
            count,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<TransactionReply>()
            .ok_or_else(|| anyhow!("invalid history transaction reply type"))?;
        Ok(reply.clone())
    }

    pub fn get_xdxr_info(&self, code: &str) -> Result<Vec<XdxrInfo>> {
        use crate::exchange::symbol::detect_market;

        let (market_id, _, symbol) = detect_market(code);
        if symbol.len() != 6 {
            return Err(anyhow!("invalid security code"));
        }

        let mut msg = new_xdxr_info()?;
        let mut code_bytes = [0u8; 6];
        code_bytes[..symbol.len()].copy_from_slice(symbol.as_bytes());

        let request = XdxrInfoRequest {
            market: market_id as u8,
            code: code_bytes,
        };
        msg.set_params(request);

        self.command(&mut msg)?;
        let reply = msg
            .reply()
            .downcast_ref::<Vec<XdxrInfo>>()
            .ok_or_else(|| anyhow!("invalid xdxr info reply type"))?;
        Ok(reply.clone())
    }
}
