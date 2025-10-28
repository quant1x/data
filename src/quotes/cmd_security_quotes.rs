use std::collections::HashMap;
use std::io::Cursor;
use std::time::{SystemTime, UNIX_EPOCH};

use anyhow::{anyhow, Result};
use byteorder::{ByteOrder, LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::exchange::{assert_index_by_market_and_code, get_market_flag};
use crate::internal::{
    base_unit, decode_varint, int_to_f64, is_nan_or_inf, sequence_id, time_from_int,
};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

pub const SECURITY_QUOTES_MAX: usize = 80;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TradeState {
    Delisting = 0,
    Normal = 1,
    Suspend = 2,
}

impl Default for TradeState {
    fn default() -> Self {
        TradeState::Normal
    }
}

#[derive(Debug, Clone)]
pub struct Stock {
    pub market: u8,
    pub code: String,
}

#[derive(Debug, Default, Clone)]
pub struct SecurityQuotesRequest {
    pub stock_list: Vec<Stock>,
}

#[derive(Debug, Clone)]
pub struct SecurityQuote {
    pub state: TradeState,
    pub market: u8,
    pub code: String,
    pub active1: u16,
    pub price: f64,
    pub last_close: f64,
    pub open: f64,
    pub high: f64,
    pub low: f64,
    pub server_time: String,
    pub reversed_bytes0: i32,
    pub reversed_bytes1: i32,
    pub vol: i32,
    pub cur_vol: i32,
    pub amount: f64,
    pub s_vol: i32,
    pub b_vol: i32,
    pub index_open_amount: i32,
    pub stock_open_amount: i32,
    pub open_volume: i32,
    pub close_volume: i32,
    pub index_up: i32,
    pub index_up_limit: i32,
    pub index_down: i32,
    pub index_down_limit: i32,
    pub bid1: f64,
    pub ask1: f64,
    pub bid_vol1: i32,
    pub ask_vol1: i32,
    pub bid2: f64,
    pub ask2: f64,
    pub bid_vol2: i32,
    pub ask_vol2: i32,
    pub bid3: f64,
    pub ask3: f64,
    pub bid_vol3: i32,
    pub ask_vol3: i32,
    pub bid4: f64,
    pub ask4: f64,
    pub bid_vol4: i32,
    pub ask_vol4: i32,
    pub bid5: f64,
    pub ask5: f64,
    pub bid_vol5: i32,
    pub ask_vol5: i32,
    pub reversed_bytes4: u16,
    pub reversed_bytes5: i32,
    pub reversed_bytes6: i32,
    pub reversed_bytes7: i32,
    pub reversed_bytes8: i32,
    pub rate: f64,
    pub active2: u16,
    pub timestamp: String,
}

impl Default for SecurityQuote {
    fn default() -> Self {
        Self {
            state: TradeState::Normal,
            market: 0,
            code: String::new(),
            active1: 0,
            price: 0.0,
            last_close: 0.0,
            open: 0.0,
            high: 0.0,
            low: 0.0,
            server_time: String::new(),
            reversed_bytes0: 0,
            reversed_bytes1: 0,
            vol: 0,
            cur_vol: 0,
            amount: 0.0,
            s_vol: 0,
            b_vol: 0,
            index_open_amount: 0,
            stock_open_amount: 0,
            open_volume: 0,
            close_volume: 0,
            index_up: 0,
            index_up_limit: 0,
            index_down: 0,
            index_down_limit: 0,
            bid1: 0.0,
            ask1: 0.0,
            bid_vol1: 0,
            ask_vol1: 0,
            bid2: 0.0,
            ask2: 0.0,
            bid_vol2: 0,
            ask_vol2: 0,
            bid3: 0.0,
            ask3: 0.0,
            bid_vol3: 0,
            ask_vol3: 0,
            bid4: 0.0,
            ask4: 0.0,
            bid_vol4: 0,
            ask_vol4: 0,
            bid5: 0.0,
            ask5: 0.0,
            bid_vol5: 0,
            ask_vol5: 0,
            reversed_bytes4: 0,
            reversed_bytes5: 0,
            reversed_bytes6: 0,
            reversed_bytes7: 0,
            reversed_bytes8: 0,
            rate: 0.0,
            active2: 0,
            timestamp: String::new(),
        }
    }
}

#[derive(Debug, Default, Clone)]
pub struct SecurityQuotesReply {
    pub count: u16,
    pub list: Vec<SecurityQuote>,
}

pub struct SecurityQuotesPackage {
    req_header: StdRequestHeader,
    request: SecurityQuotesRequest,
    resp_header: Option<StdResponseHeader>,
    reply: SecurityQuotesReply,
    map_code: HashMap<String, Stock>,
}

impl SecurityQuotesPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_SECURITY_QUOTES_OLD;
        Ok(Self {
            req_header,
            request: SecurityQuotesRequest::default(),
            resp_header: None,
            reply: SecurityQuotesReply::default(),
            map_code: HashMap::new(),
        })
    }

    pub fn set_params(&mut self, request: SecurityQuotesRequest) {
        self.request = request;
        self.map_code.clear();
        for stock in &self.request.stock_list {
            let market_flag = get_market_flag(stock.market as u16);
            self.map_code
                .insert(format!("{}{}", market_flag, stock.code), stock.clone());
        }
    }

    fn calc_payload_len(&self) -> u16 {
        let count = self.request.stock_list.len();
        2 + (count as u16 * 7) + 10
    }

    fn get_current_timestamp() -> String {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default();
        format!("{}", now.as_millis())
    }

    fn get_price(base_unit: f64, price: i32, diff: i32) -> f64 {
        (price + diff) as f64 / base_unit
    }

    fn resolve_state(open: f64, last_close: f64) -> TradeState {
        if last_close == 0.0 && open == 0.0 {
            TradeState::Delisting
        } else if open == 0.0 {
            TradeState::Suspend
        } else {
            TradeState::Normal
        }
    }
}

impl Message for SecurityQuotesPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        if self.request.stock_list.is_empty() {
            return Err(anyhow!(
                "security quotes request requires at least one stock"
            ));
        }
        if self.request.stock_list.len() > SECURITY_QUOTES_MAX {
            return Err(anyhow!(
                "security quotes request exceeds limit {}",
                SECURITY_QUOTES_MAX
            ));
        }

        self.req_header.seq_id = sequence_id();
        let payload_len = self.calc_payload_len();
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        // 历史兼容字段，拷贝 Go 版本常量。
        out.extend_from_slice(&[0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00]);
        out.write_u16::<LittleEndian>(self.request.stock_list.len() as u16)?;
        for stock in &self.request.stock_list {
            out.write_u8(stock.market)?;
            let mut code_bytes = [0u8; 6];
            let len = stock.code.as_bytes().len().min(6);
            code_bytes[..len].copy_from_slice(&stock.code.as_bytes()[..len]);
            out.extend_from_slice(&code_bytes);
        }
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);
        if cursor.get_ref().len() < 4 {
            return Err(anyhow!("security quotes response too short"));
        }
        cursor.read_u16::<LittleEndian>()?; // 跳过两个字节
        self.reply.count = cursor.read_u16::<LittleEndian>()?;
        self.reply.list.clear();
        self.reply.list.reserve(self.reply.count as usize);

        let data = cursor.get_ref();
        let mut pos = cursor.position() as usize;

        for _ in 0..self.reply.count {
            let mut entry = SecurityQuote::default();
            if pos >= data.len() {
                return Err(anyhow!("security quotes entry truncated"));
            }
            entry.market = data[pos];
            pos += 1;
            let mut code_bytes = [0u8; 6];
            if pos + 6 > data.len() {
                return Err(anyhow!("security quotes code overflow"));
            }
            code_bytes.copy_from_slice(&data[pos..pos + 6]);
            pos += 6;
            entry.code = String::from_utf8_lossy(&code_bytes)
                .trim_end_matches('\u{0}')
                .to_string();

            if pos + 2 > data.len() {
                return Err(anyhow!("security quotes active1 overflow"));
            }
            entry.active1 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;

            let price = decode_varint(data, &mut pos)?;
            let base_unit = base_unit(entry.market as u16, &entry.code);
            entry.price = Self::get_price(base_unit, price, 0);
            entry.last_close = Self::get_price(base_unit, price, decode_varint(data, &mut pos)?);
            entry.open = Self::get_price(base_unit, price, decode_varint(data, &mut pos)?);
            entry.high = Self::get_price(base_unit, price, decode_varint(data, &mut pos)?);
            entry.low = Self::get_price(base_unit, price, decode_varint(data, &mut pos)?);

            entry.reversed_bytes0 = decode_varint(data, &mut pos)?;
            if entry.reversed_bytes0 > 0 {
                entry.server_time = time_from_int(entry.reversed_bytes0);
            } else {
                entry.server_time = "0".to_string();
            }
            entry.reversed_bytes1 = decode_varint(data, &mut pos)?;
            entry.vol = decode_varint(data, &mut pos)? * 100;
            entry.cur_vol = decode_varint(data, &mut pos)?;

            if pos + 4 > data.len() {
                return Err(anyhow!("security quotes amount overflow"));
            }
            let amount_raw = LittleEndian::read_u32(&data[pos..pos + 4]);
            pos += 4;
            entry.amount = int_to_f64(amount_raw as i32);

            entry.s_vol = decode_varint(data, &mut pos)?;
            entry.b_vol = decode_varint(data, &mut pos)?;
            entry.index_open_amount = decode_varint(data, &mut pos)? * 100;
            entry.stock_open_amount = decode_varint(data, &mut pos)? * 100;

            let is_index_or_block =
                assert_index_by_market_and_code(entry.market as u16, &entry.code);
            let mut open_volume = if entry.open == 0.0 {
                0.0
            } else if is_index_or_block {
                (entry.index_open_amount as f64 / entry.open).round()
            } else {
                (entry.stock_open_amount as f64 / entry.open).round()
            };
            if is_nan_or_inf(open_volume) {
                open_volume = 0.0;
            }
            entry.open_volume = open_volume as i32;

            let mut bid_prices = [0f64; 5];
            let mut ask_prices = [0f64; 5];
            let mut bid_vols = [0i32; 5];
            let mut ask_vols = [0i32; 5];
            for i in 0..5 {
                bid_prices[i] = Self::get_price(base_unit, decode_varint(data, &mut pos)?, price);
                ask_prices[i] = Self::get_price(base_unit, decode_varint(data, &mut pos)?, price);
                bid_vols[i] = decode_varint(data, &mut pos)?;
                ask_vols[i] = decode_varint(data, &mut pos)?;
            }

            entry.bid1 = bid_prices[0];
            entry.bid2 = bid_prices[1];
            entry.bid3 = bid_prices[2];
            entry.bid4 = bid_prices[3];
            entry.bid5 = bid_prices[4];
            entry.ask1 = ask_prices[0];
            entry.ask2 = ask_prices[1];
            entry.ask3 = ask_prices[2];
            entry.ask4 = ask_prices[3];
            entry.ask5 = ask_prices[4];

            entry.bid_vol1 = bid_vols[0];
            entry.bid_vol2 = bid_vols[1];
            entry.bid_vol3 = bid_vols[2];
            entry.bid_vol4 = bid_vols[3];
            entry.bid_vol5 = bid_vols[4];
            entry.ask_vol1 = ask_vols[0];
            entry.ask_vol2 = ask_vols[1];
            entry.ask_vol3 = ask_vols[2];
            entry.ask_vol4 = ask_vols[3];
            entry.ask_vol5 = ask_vols[4];

            if is_index_or_block {
                entry.index_up = entry.bid_vol1;
                entry.index_down = entry.ask_vol1;
                entry.index_up_limit = entry.bid_vol2;
                entry.index_down_limit = entry.ask_vol2;
            }

            if pos + 2 > data.len() {
                return Err(anyhow!("security quotes reversed bytes4 overflow"));
            }
            entry.reversed_bytes4 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;
            entry.reversed_bytes5 = decode_varint(data, &mut pos)?;
            entry.reversed_bytes6 = decode_varint(data, &mut pos)?;
            entry.reversed_bytes7 = decode_varint(data, &mut pos)?;
            entry.reversed_bytes8 = decode_varint(data, &mut pos)?;

            if pos + 2 > data.len() {
                return Err(anyhow!("security quotes rate overflow"));
            }
            let reversed_bytes9 = LittleEndian::read_i16(&data[pos..pos + 2]);
            pos += 2;
            entry.rate = reversed_bytes9 as f64 / 100.0;

            // 跳过保留 2 字节
            pos += 2;
            if pos + 2 > data.len() {
                return Err(anyhow!("security quotes active2 overflow"));
            }
            entry.active2 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;

            entry.state = Self::resolve_state(entry.open, entry.last_close);
            if entry.state != TradeState::Delisting {
                let market_flag = get_market_flag(entry.market as u16);
                let map_key = format!("{}{}", market_flag, entry.code);
                self.map_code.remove(&map_key);
            }

            entry.timestamp = Self::get_current_timestamp();
            self.reply.list.push(entry);
        }

        if !self.map_code.is_empty() {
            for quote in self.reply.list.iter_mut() {
                if self.map_code.is_empty() {
                    break;
                }

                if quote.state != TradeState::Delisting {
                    continue;
                }

                let market_flag = get_market_flag(quote.market as u16);
                let key = format!("{}{}", market_flag, quote.code);
                if self.map_code.remove(&key).is_some() {
                    continue;
                }

                let candidate = self.map_code.iter().find_map(|(k, stock)| {
                    if stock.market == quote.market {
                        Some((k.clone(), stock.code.clone()))
                    } else {
                        None
                    }
                });

                if let Some((k, corrected_code)) = candidate {
                    quote.code = corrected_code;
                    self.map_code.remove(&k);
                }
            }
        }

        self.map_code.clear();

        cursor.set_position(pos as u64);
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_security_quotes() -> Result<SecurityQuotesPackage> {
    SecurityQuotesPackage::new()
}
