use crate::internal::{decode_varint, get_datetime, int_to_f64, sequence_id};
use crate::proto::{FLAG_NOT_ZIPPED, STD_MSG_SECURITY_BARS};
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};
use anyhow::Result;
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::io::Cursor;

pub const SECURITY_BARS_MAX: usize = 800;

#[derive(Debug, Clone)]
pub struct SecurityBarsRequest {
    pub market: u16,
    pub code: [u8; 6],
    pub category: u16, // 种类 5分钟  10分钟
    pub i: u16,        // 未知 填充, 间隔多少个Category
    pub start: u16,
    pub count: u16,
}

#[derive(Debug, Clone)]
pub struct SecurityBarsReply {
    pub count: u16,
    pub list: Vec<SecurityBar>,
}

#[derive(Debug, Clone)]
pub struct SecurityBar {
    pub open: f64,
    pub close: f64,
    pub high: f64,
    pub low: f64,
    pub vol: f64,
    pub amount: f64,
    pub year: i32,
    pub month: i32,
    pub day: i32,
    pub hour: i32,
    pub minute: i32,
    pub date_time: String,
    pub up_count: u16,   // 指数有效, 上涨家数
    pub down_count: u16, // 指数有效, 下跌家数
}

impl Default for SecurityBar {
    fn default() -> Self {
        SecurityBar {
            open: 0.0,
            close: 0.0,
            high: 0.0,
            low: 0.0,
            vol: 0.0,
            amount: 0.0,
            year: 0,
            month: 0,
            day: 0,
            hour: 0,
            minute: 0,
            date_time: String::new(),
            up_count: 0,
            down_count: 0,
        }
    }
}

pub struct SecurityBarsPackage {
    req_header: StdRequestHeader,
    request: SecurityBarsRequest,
    resp_header: Option<StdResponseHeader>,
    reply: SecurityBarsReply,
}

impl SecurityBarsPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x00;
        req_header.method = STD_MSG_SECURITY_BARS;

        Ok(Self {
            req_header,
            request: SecurityBarsRequest {
                market: 0,
                code: [0; 6],
                category: 0,
                i: 0,
                start: 0,
                count: 0,
            },
            resp_header: None,
            reply: SecurityBarsReply {
                count: 0,
                list: Vec::new(),
            },
        })
    }

    pub fn set_params(&mut self, request: SecurityBarsRequest) {
        self.request = request;
        if self.request.i < 1 {
            self.request.i = 1;
        }
    }
}

impl Message for SecurityBarsPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x1c;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u16::<LittleEndian>(self.request.category)?;
        out.write_u16::<LittleEndian>(self.request.i)?;
        out.write_u16::<LittleEndian>(self.request.start)?;
        out.write_u16::<LittleEndian>(self.request.count)?;

        // Add fixed content (20 bytes of zeros)
        out.extend_from_slice(&[0u8; 20]);

        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        self.resp_header = Some(*header);

        let mut cursor = Cursor::new(body);
        let count = cursor.read_u16::<LittleEndian>()?;
        self.reply.count = count;

        let mut pos = 2;
        let mut pre_diff_base = 0i32;

        for _ in 0..count {
            let mut bar = SecurityBar::default();

            // Get datetime
            let (year, month, day, hour, minute) =
                get_datetime(self.request.category as i32, body, &mut pos)?;
            bar.year = year;
            bar.month = month;
            bar.day = day;
            bar.hour = hour;
            bar.minute = minute;
            bar.date_time = format!(
                "{:04}-{:02}-{:02} {:02}:{:02}:00.000",
                year, month, day, hour, minute
            );

            // Decode price diffs
            let price_open_diff = decode_varint(body, &mut pos)?;
            let price_close_diff = decode_varint(body, &mut pos)?;
            let price_high_diff = decode_varint(body, &mut pos)?;
            let price_low_diff = decode_varint(body, &mut pos)?;

            // Read volume and amount
            let ivol = cursor.read_u32::<LittleEndian>()?;
            bar.vol = int_to_f64(ivol);
            pos += 4;

            let dbvol = cursor.read_u32::<LittleEndian>()?;
            bar.amount = int_to_f64(dbvol as i32);
            pos += 4;

            // For index bars, read up/down counts (but this is security bars, so skip for now)
            // These would be read in IndexBarsPackage

            // Calculate actual prices
            bar.open = (price_open_diff + pre_diff_base) as f64 / 1000.0;
            let price_open = price_open_diff + pre_diff_base;

            bar.close = (price_open + price_close_diff) as f64 / 1000.0;
            bar.high = (price_open + price_high_diff) as f64 / 1000.0;
            bar.low = (price_open + price_low_diff) as f64 / 1000.0;

            pre_diff_base = price_open + price_close_diff;

            self.reply.list.push(bar);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_security_bars() -> Result<SecurityBarsPackage> {
    SecurityBarsPackage::new()
}

pub struct IndexBarsPackage {
    req_header: StdRequestHeader,
    request: SecurityBarsRequest,
    resp_header: Option<StdResponseHeader>,
    reply: SecurityBarsReply,
}

impl IndexBarsPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x00;
        req_header.method = crate::proto::STD_MSG_INDEXBARS;

        Ok(Self {
            req_header,
            request: SecurityBarsRequest {
                market: 0,
                code: [0; 6],
                category: 0,
                i: 0,
                start: 0,
                count: 0,
            },
            resp_header: None,
            reply: SecurityBarsReply {
                count: 0,
                list: Vec::new(),
            },
        })
    }

    pub fn set_params(&mut self, request: SecurityBarsRequest) {
        self.request = request;
        self.request.i = 1; // Index bars always set I to 1
    }
}

impl Message for IndexBarsPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x1c;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u16::<LittleEndian>(self.request.category)?;
        out.write_u16::<LittleEndian>(self.request.i)?;
        out.write_u16::<LittleEndian>(self.request.start)?;
        out.write_u16::<LittleEndian>(self.request.count)?;

        // Add fixed content (20 bytes of zeros)
        out.extend_from_slice(&[0u8; 20]);

        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        self.resp_header = Some(*header);

        let mut cursor = Cursor::new(body);
        let count = cursor.read_u16::<LittleEndian>()?;
        self.reply.count = count;

        let mut pos = 2;
        let mut pre_diff_base = 0i32;

        for _ in 0..count {
            let mut bar = SecurityBar::default();

            // Get datetime
            let (year, month, day, hour, minute) =
                get_datetime(self.request.category as i32, body, &mut pos)?;
            bar.year = year;
            bar.month = month;
            bar.day = day;
            bar.hour = hour;
            bar.minute = minute;
            bar.date_time = format!(
                "{:04}-{:02}-{:02} {:02}:{:02}:00.000",
                year, month, day, hour, minute
            );

            // Decode price diffs
            let price_open_diff = decode_varint(body, &mut pos)?;
            let price_close_diff = decode_varint(body, &mut pos)?;
            let price_high_diff = decode_varint(body, &mut pos)?;
            let price_low_diff = decode_varint(body, &mut pos)?;

            // Read volume and amount
            let ivol = cursor.read_u32::<LittleEndian>()?;
            bar.vol = int_to_f64(ivol);
            pos += 4;

            let dbvol = cursor.read_u32::<LittleEndian>()?;
            bar.amount = int_to_f64(dbvol as i32);
            pos += 4;

            // For index bars, read up/down counts
            bar.up_count = cursor.read_u16::<LittleEndian>()?;
            pos += 2;
            bar.down_count = cursor.read_u16::<LittleEndian>()?;
            pos += 2;

            // Calculate actual prices
            bar.open = (price_open_diff + pre_diff_base) as f64 / 1000.0;
            let price_open = price_open_diff + pre_diff_base;

            bar.close = (price_open + price_close_diff) as f64 / 1000.0;
            bar.high = (price_open + price_high_diff) as f64 / 1000.0;
            bar.low = (price_open + price_low_diff) as f64 / 1000.0;

            pre_diff_base = price_open + price_close_diff;

            self.reply.list.push(bar);
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_index_bars() -> Result<IndexBarsPackage> {
    IndexBarsPackage::new()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_security_bar_default() {
        let bar = SecurityBar::default();
        assert_eq!(bar.open, 0.0);
        assert_eq!(bar.close, 0.0);
        assert_eq!(bar.high, 0.0);
        assert_eq!(bar.low, 0.0);
        assert_eq!(bar.vol, 0.0);
        assert_eq!(bar.amount, 0.0);
        assert_eq!(bar.year, 0);
        assert_eq!(bar.month, 0);
        assert_eq!(bar.day, 0);
        assert_eq!(bar.hour, 0);
        assert_eq!(bar.minute, 0);
        assert_eq!(bar.date_time, "");
        assert_eq!(bar.up_count, 0);
        assert_eq!(bar.down_count, 0);
    }

    #[test]
    fn test_security_bars_package_new() {
        let package = new_security_bars().unwrap();
        assert_eq!(package.request.market, 0);
        assert_eq!(package.request.count, 0);
        assert_eq!(package.reply.count, 0);
        assert!(package.reply.list.is_empty());
    }
}