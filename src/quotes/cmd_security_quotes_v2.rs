use std::io::Cursor;

use anyhow::{anyhow, Result};
use byteorder::{ByteOrder, LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::internal::{decode_varint, int_to_f64, sequence_id, time_from_int, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

pub const SECURITY_QUOTES_MAX_V2: usize = 80;

#[derive(Debug, Clone)]
pub struct V2Stock {
    pub market: u8,
    pub code: String,
}

#[derive(Debug, Default, Clone)]
pub struct V2SecurityQuotesRequest {
    pub stock_list: Vec<V2Stock>,
}

#[derive(Debug, Clone)]
pub struct V2SecurityQuote {
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
    pub reversed_bytes2: i32,
    pub reversed_bytes3: i32,
    pub bid1: f64,
    pub ask1: f64,
    pub bid_vol1: i32,
    pub ask_vol1: i32,
    pub reversed_bytes4: u16,
    pub reversed_bytes5: i32,
    pub reversed_bytes6: i32,
    pub reversed_bytes7: i32,
    pub reversed_bytes8: i32,
    pub rate: f64,
    pub active2: u16,
}

impl Default for V2SecurityQuote {
    fn default() -> Self {
        Self {
            market: 0,
            code: String::new(),
            active1: 0,
            price: 0.0,
            last_close: 0.0,
            open: 0.0,
            high: 0.0,
            low: 0.0,
            server_time: "0".to_string(),
            reversed_bytes0: 0,
            reversed_bytes1: 0,
            vol: 0,
            cur_vol: 0,
            amount: 0.0,
            s_vol: 0,
            b_vol: 0,
            reversed_bytes2: 0,
            reversed_bytes3: 0,
            bid1: 0.0,
            ask1: 0.0,
            bid_vol1: 0,
            ask_vol1: 0,
            reversed_bytes4: 0,
            reversed_bytes5: 0,
            reversed_bytes6: 0,
            reversed_bytes7: 0,
            reversed_bytes8: 0,
            rate: 0.0,
            active2: 0,
        }
    }
}

#[derive(Debug, Default, Clone)]
pub struct V2SecurityQuotesReply {
    pub count: u16,
    pub list: Vec<V2SecurityQuote>,
}

pub struct V2SecurityQuotesPackage {
    req_header: StdRequestHeader,
    request: V2SecurityQuotesRequest,
    resp_header: Option<StdResponseHeader>,
    reply: V2SecurityQuotesReply,
}

impl V2SecurityQuotesPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_SECURITY_QUOTES_NEW;
        Ok(Self {
            req_header,
            request: V2SecurityQuotesRequest::default(),
            resp_header: None,
            reply: V2SecurityQuotesReply::default(),
        })
    }

    pub fn set_params(&mut self, request: V2SecurityQuotesRequest) {
        self.request = request;
    }

    fn calc_payload_len(&self) -> u16 {
        let count = self.request.stock_list.len();
        2 + (count as u16 * 7) + 10
    }
}

impl Message for V2SecurityQuotesPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        if self.request.stock_list.is_empty() {
            return Err(anyhow!(
                "v2 security quotes request requires at least one stock"
            ));
        }
        if self.request.stock_list.len() > SECURITY_QUOTES_MAX_V2 {
            return Err(anyhow!(
                "v2 security quotes request exceeds limit {}",
                SECURITY_QUOTES_MAX_V2
            ));
        }

        self.req_header.seq_id = sequence_id();
        let payload_len = self.calc_payload_len();
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
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
        if body.len() < 4 {
            return Err(anyhow!("v2 security quotes response too short"));
        }
        let mut cursor = Cursor::new(body);
        cursor.read_u16::<LittleEndian>()?; // 跳过两个字节
        self.reply.count = cursor.read_u16::<LittleEndian>()?;
        let data = cursor.get_ref();
        let mut pos = cursor.position() as usize;
        self.reply.list.clear();
        self.reply.list.reserve(self.reply.count as usize);

        for _ in 0..self.reply.count {
            let mut entry = V2SecurityQuote::default();
            if pos >= data.len() {
                return Err(anyhow!("v2 security quotes entry truncated"));
            }
            entry.market = data[pos];
            pos += 1;

            if pos + 6 > data.len() {
                return Err(anyhow!("v2 security quotes code overflow"));
            }
            let code_bytes = &data[pos..pos + 6];
            pos += 6;
            entry.code = utf8_to_gbk(code_bytes);

            if pos + 2 > data.len() {
                return Err(anyhow!("v2 security quotes active1 overflow"));
            }
            entry.active1 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;

            let price = decode_varint(data, &mut pos)?;
            entry.price = price as f64 / 100.0;
            entry.last_close = (price + decode_varint(data, &mut pos)?) as f64 / 100.0;
            entry.open = (price + decode_varint(data, &mut pos)?) as f64 / 100.0;
            entry.high = (price + decode_varint(data, &mut pos)?) as f64 / 100.0;
            entry.low = (price + decode_varint(data, &mut pos)?) as f64 / 100.0;

            entry.reversed_bytes0 = decode_varint(data, &mut pos)?;
            if entry.reversed_bytes0 > 0 {
                entry.server_time = time_from_int(entry.reversed_bytes0);
            }
            entry.reversed_bytes1 = decode_varint(data, &mut pos)?;
            entry.vol = decode_varint(data, &mut pos)?;
            entry.cur_vol = decode_varint(data, &mut pos)?;

            if pos + 4 > data.len() {
                return Err(anyhow!("v2 security quotes amount overflow"));
            }
            let amount_raw = LittleEndian::read_u32(&data[pos..pos + 4]);
            pos += 4;
            entry.amount = int_to_f64(amount_raw as i32);

            entry.s_vol = decode_varint(data, &mut pos)?;
            entry.b_vol = decode_varint(data, &mut pos)?;
            entry.reversed_bytes2 = decode_varint(data, &mut pos)?;
            entry.reversed_bytes3 = decode_varint(data, &mut pos)?;

            let bid_price = decode_varint(data, &mut pos)?;
            let ask_price = decode_varint(data, &mut pos)?;
            entry.bid1 = (bid_price + price) as f64 / 100.0;
            entry.ask1 = (ask_price + price) as f64 / 100.0;
            entry.bid_vol1 = decode_varint(data, &mut pos)?;
            entry.ask_vol1 = decode_varint(data, &mut pos)?;

            if pos + 2 > data.len() {
                return Err(anyhow!("v2 security quotes reversed bytes4 overflow"));
            }
            entry.reversed_bytes4 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;

            // 跳过扩展字段，保持与 Go 逻辑一致
            for _ in 0..2 {
                decode_varint(data, &mut pos)?;
            }
            let skip_len = 12 * 4;
            if pos + skip_len > data.len() {
                return Err(anyhow!("v2 security quotes extended fields overflow"));
            }
            pos += skip_len;

            if pos + 2 > data.len() {
                return Err(anyhow!("v2 security quotes active2 overflow"));
            }
            entry.active2 = LittleEndian::read_u16(&data[pos..pos + 2]);
            pos += 2;

            self.reply.list.push(entry);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_security_quotes_v2() -> Result<V2SecurityQuotesPackage> {
    V2SecurityQuotesPackage::new()
}
