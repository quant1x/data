use crate::internal::helpers::{base_unit, decode_varint, get_time, sequence_id};
use crate::proto::{STD_MSG_HISTORY_TRANSACTION_DATA, STD_MSG_TRANSACTION_DATA};
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};
use anyhow::Result;
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::any::Any;
use std::io::Cursor;

pub const TRANSACTION_MAX: usize = 1800;

#[derive(Debug, Clone, PartialEq)]
pub enum TradeType {
    Buy = 0,     // 买入
    Sell = 1,    // 卖出
    Neutral = 2, // 中性盘
    Unknown = 3, // 未知类型
}

#[derive(Debug, Clone, Default)]
pub struct TickTransaction {
    pub time: String,
    pub price: f64,
    pub vol: i32,
    pub num: i32,
    pub amount: f64,
    pub buy_or_sell: i32,
}

#[derive(Debug, Clone, Default)]
pub struct TransactionRequest {
    pub market: u16,
    pub code: [u8; 6],
    pub start: u16,
    pub count: u16,
}

#[derive(Debug, Clone, Default)]
pub struct HistoryTransactionRequest {
    pub date: u32,
    pub market: u16,
    pub code: [u8; 6],
    pub start: u16,
    pub count: u16,
}

#[derive(Debug, Clone, Default)]
pub struct TransactionReply {
    pub count: u16,
    pub list: Vec<TickTransaction>,
}

#[derive(Debug, Clone)]
pub struct TransactionPackage {
    req_header: StdRequestHeader,
    resp_header: StdResponseHeader,
    request: TransactionRequest,
    reply: TransactionReply,
}

impl TransactionPackage {
    pub fn new() -> Self {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = 0x0c; // FlagNotZipped
        req_header.seq_id = sequence_id();
        req_header.packet_type = 0x00;
        req_header.method = STD_MSG_TRANSACTION_DATA;

        Self {
            req_header,
            resp_header: StdResponseHeader::default(),
            request: TransactionRequest::default(),
            reply: TransactionReply::default(),
        }
    }

    pub fn set_params(&mut self, req: TransactionRequest) {
        self.request = req;
    }
}

impl Message for TransactionPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x0e;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u16::<LittleEndian>(self.request.start)?;
        out.write_u16::<LittleEndian>(self.request.count)?;

        Ok(out)
    }

    fn unserialize(&mut self, _header: &StdResponseHeader, data: &[u8]) -> Result<()> {
        self.resp_header = _header.clone();

        let mut cursor = Cursor::new(data);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;

        let market = crate::exchange::MarketType::from(self.request.market);
        let code = String::from_utf8_lossy(&self.request.code)
            .trim_end_matches('\0')
            .to_string();
        let base_unit = base_unit(market as u16, &code);
        let is_index = crate::exchange::assert_index_by_market_and_code(market, &code);

        let mut last_price = 0i32;
        self.reply.list = Vec::with_capacity(self.reply.count as usize);

        let mut pos = cursor.position() as usize;
        for _ in 0..self.reply.count {
            let mut transaction = TickTransaction::default();

            let (hour, minute) = get_time(data, &mut pos)?;
            transaction.time = format!("{:02}:{:02}", hour, minute);

            let raw_price = decode_varint(data, &mut pos)?;
            let vol = decode_varint(data, &mut pos)?;
            let num = decode_varint(data, &mut pos)?;
            let buy_or_sell = decode_varint(data, &mut pos)?;

            last_price += raw_price;
            transaction.price = last_price as f64 / base_unit;
            transaction.num = num as i32;
            transaction.buy_or_sell = buy_or_sell as i32;

            if is_index {
                let amount = vol * 100;
                transaction.amount = amount as f64;
                transaction.vol = ((amount as f64) / transaction.price) as i32;
            } else {
                transaction.vol = (vol * 100) as i32;
                transaction.amount = transaction.vol as f64 * transaction.price;
            }

            // Skip one varint (unused)
            decode_varint(data, &mut pos)?;

            self.reply.list.push(transaction);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn Any + Send + Sync) {
        &self.reply
    }
}

#[derive(Debug, Clone)]
pub struct HistoryTransactionPackage {
    req_header: StdRequestHeader,
    resp_header: StdResponseHeader,
    request: HistoryTransactionRequest,
    reply: TransactionReply,
}

impl HistoryTransactionPackage {
    pub fn new() -> Self {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = 0x0c; // FlagNotZipped
        req_header.seq_id = sequence_id();
        req_header.packet_type = 0x00;
        req_header.method = STD_MSG_HISTORY_TRANSACTION_DATA;

        Self {
            req_header,
            resp_header: StdResponseHeader::default(),
            request: HistoryTransactionRequest::default(),
            reply: TransactionReply::default(),
        }
    }

    pub fn set_params(&mut self, req: HistoryTransactionRequest) {
        self.request = req;
    }
}

impl Message for HistoryTransactionPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x12;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u32::<LittleEndian>(self.request.date)?;
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u16::<LittleEndian>(self.request.start)?;
        out.write_u16::<LittleEndian>(self.request.count)?;

        Ok(out)
    }

    fn unserialize(&mut self, _header: &StdResponseHeader, data: &[u8]) -> Result<()> {
        self.resp_header = _header.clone();

        let mut cursor = Cursor::new(data);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;

        // Skip 4 bytes
        cursor.set_position(cursor.position() + 4);

        let market = crate::exchange::MarketType::from(self.request.market);
        let code = String::from_utf8_lossy(&self.request.code)
            .trim_end_matches('\0')
            .to_string();
        let base_unit = base_unit(market as u16, &code);
        let is_index = crate::exchange::assert_index_by_market_and_code(market, &code);

        let mut last_price = 0i32;
        self.reply.list = Vec::with_capacity(self.reply.count as usize);

        let mut pos = cursor.position() as usize;
        for _ in 0..self.reply.count {
            let mut transaction = TickTransaction::default();

            let (hour, minute) = get_time(data, &mut pos)?;
            transaction.time = format!("{:02}:{:02}", hour, minute);

            let raw_price = decode_varint(data, &mut pos)?;
            let vol = decode_varint(data, &mut pos)?;
            let buy_or_sell = decode_varint(data, &mut pos)?;

            // Skip one varint (num field not present in history data)
            decode_varint(data, &mut pos)?;

            last_price += raw_price;
            transaction.price = last_price as f64 / base_unit;
            transaction.buy_or_sell = buy_or_sell as i32;

            if is_index {
                let amount = vol * 100;
                transaction.amount = amount as f64;
                transaction.vol = ((amount as f64) / transaction.price) as i32;
            } else {
                transaction.vol = (vol * 100) as i32;
                transaction.amount = transaction.vol as f64 * transaction.price;
            }

            self.reply.list.push(transaction);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_transaction() -> Result<TransactionPackage> {
    Ok(TransactionPackage::new())
}

pub fn new_history_transaction() -> Result<HistoryTransactionPackage> {
    Ok(HistoryTransactionPackage::new())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tick_transaction_default() {
        let tt = TickTransaction::default();
        assert_eq!(tt.time, "");
        assert_eq!(tt.price, 0.0);
        assert_eq!(tt.vol, 0);
        assert_eq!(tt.num, 0);
        assert_eq!(tt.amount, 0.0);
        assert_eq!(tt.buy_or_sell, 0);
    }

    #[test]
    fn test_transaction_package_new() {
        let pkg = TransactionPackage::new();
        assert_eq!(pkg.req_header.method, STD_MSG_TRANSACTION_DATA);
        assert_eq!(pkg.req_header.zip_flag, 0x0c);
    }

    #[test]
    fn test_history_transaction_package_new() {
        let pkg = HistoryTransactionPackage::new();
        assert_eq!(pkg.req_header.method, STD_MSG_HISTORY_TRANSACTION_DATA);
        assert_eq!(pkg.req_header.zip_flag, 0x0c);
    }

    #[test]
    fn test_trade_type_values() {
        assert_eq!(TradeType::Buy as i32, 0);
        assert_eq!(TradeType::Sell as i32, 1);
        assert_eq!(TradeType::Neutral as i32, 2);
        assert_eq!(TradeType::Unknown as i32, 3);
    }
}