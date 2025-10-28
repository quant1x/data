use crate::internal::{base_unit, decode_varint, sequence_id};
use crate::proto::{STD_MSG_HISTORY_MINUTETIME_DATA, STD_MSG_MINUTETIME_DATA};
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};
use anyhow::Result;
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::any::Any;
use std::io::Cursor;

pub const MINUTE_TIME_MAX: usize = 1000;

#[derive(Debug, Clone, Default)]
pub struct MinuteTime {
    pub price: f32,
    pub vol: i32,
}

#[derive(Debug, Clone, Default)]
pub struct MinuteTimeRequest {
    pub market: u16,
    pub code: [u8; 6],
    pub date: u32,
}

#[derive(Debug, Clone, Default)]
pub struct HistoryMinuteTimeRequest {
    pub date: u32,
    pub market: u8,
    pub code: [u8; 6],
}

#[derive(Debug, Clone, Default)]
pub struct MinuteTimeReply {
    pub count: u16,
    pub list: Vec<MinuteTime>,
}

#[derive(Debug, Clone)]
pub struct MinuteTimePackage {
    req_header: StdRequestHeader,
    resp_header: StdResponseHeader,
    request: MinuteTimeRequest,
    reply: MinuteTimeReply,
}

impl MinuteTimePackage {
    pub fn new() -> Self {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = 0x0c; // FlagNotZipped
        req_header.seq_id = sequence_id();
        req_header.packet_type = 0x00;
        req_header.method = STD_MSG_MINUTETIME_DATA;

        Self {
            req_header,
            resp_header: StdResponseHeader::default(),
            request: MinuteTimeRequest::default(),
            reply: MinuteTimeReply::default(),
        }
    }

    pub fn set_params(&mut self, req: MinuteTimeRequest) {
        self.request = req;
    }
}

impl Message for MinuteTimePackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x0e;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u32::<LittleEndian>(self.request.date)?;

        Ok(out)
    }

    fn unserialize(&mut self, _header: &StdResponseHeader, data: &[u8]) -> Result<()> {
        self.resp_header = _header.clone();

        let mut cursor = Cursor::new(data);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;

        // Skip 6 bytes (4 + 2 more)
        cursor.set_position(cursor.position() + 6);

        // Skip 3 more bytes
        cursor.set_position(cursor.position() + 3);

        let market = crate::exchange::MarketType::from(self.request.market);
        let code = String::from_utf8_lossy(&self.request.code)
            .trim_end_matches('\0')
            .to_string();
        let base_unit = base_unit(market as u16, &code);

        let mut last_price = 0i32;
        self.reply.list = Vec::with_capacity(self.reply.count as usize);

        let mut pos = 11; // After count (2) + skip bytes (6 + 3)
        for _ in 0..self.reply.count {
            let raw_price = decode_varint(data, &mut pos)?;
            let reversed1 = decode_varint(data, &mut pos)?;
            let _ = reversed1; // unused
            let vol = decode_varint(data, &mut pos)?;

            last_price += raw_price;
            let price = last_price as f32 / base_unit as f32;

            self.reply.list.push(MinuteTime {
                price,
                vol: vol as i32,
            });
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn Any + Send + Sync) {
        &self.reply
    }
}

impl MinuteTimeRequest {
    pub fn serialize(&self) -> Result<Vec<u8>> {
        let mut buf = Vec::new();
        buf.extend_from_slice(&self.market.to_le_bytes());
        buf.extend_from_slice(&self.code);
        buf.extend_from_slice(&self.date.to_le_bytes());
        Ok(buf)
    }
}

#[derive(Debug, Clone)]
pub struct HistoryMinuteTimePackage {
    req_header: StdRequestHeader,
    resp_header: StdResponseHeader,
    request: HistoryMinuteTimeRequest,
    reply: MinuteTimeReply,
}

impl HistoryMinuteTimePackage {
    pub fn new() -> Self {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = 0x0c; // FlagNotZipped
        req_header.seq_id = sequence_id();
        req_header.packet_type = 0x00;
        req_header.method = STD_MSG_HISTORY_MINUTETIME_DATA;

        Self {
            req_header,
            resp_header: StdResponseHeader::default(),
            request: HistoryMinuteTimeRequest::default(),
            reply: MinuteTimeReply::default(),
        }
    }

    pub fn set_params(&mut self, req: HistoryMinuteTimeRequest) {
        self.request = req;
    }
}

impl Message for HistoryMinuteTimePackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 0x0d;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;

        // Write request
        out.write_u32::<LittleEndian>(self.request.date)?;
        out.write_u8(self.request.market)?;
        out.extend_from_slice(&self.request.code);

        Ok(out)
    }

    fn unserialize(&mut self, _header: &StdResponseHeader, data: &[u8]) -> Result<()> {
        self.resp_header = _header.clone();

        if data.len() < 2 {
            return Ok(());
        }

        let mut cursor = Cursor::new(data);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;

        if self.reply.count == 0 {
            return Ok(());
        }

        if data.len() < 6 {
            return Ok(());
        }

        // Skip 4 bytes (unknown functionality)
        cursor.set_position(cursor.position() + 4);

        let market = crate::exchange::MarketType::from(self.request.market as u16);
        let code = String::from_utf8_lossy(&self.request.code)
            .trim_end_matches('\0')
            .to_string();
        let base_unit = base_unit(market as u16, &code);

        let mut last_price = 0i32;
        self.reply.list = Vec::with_capacity(self.reply.count as usize);

        let mut pos = 6; // After count (2) + skip bytes (4)
        for _ in 0..self.reply.count {
            let raw_price = decode_varint(data, &mut pos)?;
            let reversed1 = decode_varint(data, &mut pos)?;
            let _ = reversed1; // unused
            let vol = decode_varint(data, &mut pos)?;

            last_price += raw_price;
            let price = last_price as f32 / base_unit as f32;

            self.reply.list.push(MinuteTime {
                price,
                vol: vol as i32,
            });
        }

        Ok(())
    }

    fn reply(&self) -> &(dyn Any + Send + Sync) {
        &self.reply
    }
}

impl HistoryMinuteTimeRequest {
    pub fn serialize(&self) -> Result<Vec<u8>> {
        let mut buf = Vec::new();
        buf.extend_from_slice(&self.date.to_le_bytes());
        buf.push(self.market);
        buf.extend_from_slice(&self.code);
        Ok(buf)
    }
}

pub fn new_minute_time() -> Result<MinuteTimePackage> {
    Ok(MinuteTimePackage::new())
}

pub fn new_history_minute_time() -> Result<HistoryMinuteTimePackage> {
    Ok(HistoryMinuteTimePackage::new())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_minute_time_default() {
        let mt = MinuteTime::default();
        assert_eq!(mt.price, 0.0);
        assert_eq!(mt.vol, 0);
    }

    #[test]
    fn test_minute_time_package_new() {
        let pkg = MinuteTimePackage::new();
        assert_eq!(pkg.req_header.method, STD_MSG_MINUTETIME_DATA);
        assert_eq!(pkg.req_header.zip_flag, 0x0c);
    }

    #[test]
    fn test_history_minute_time_package_new() {
        let pkg = HistoryMinuteTimePackage::new();
        assert_eq!(pkg.req_header.method, STD_MSG_HISTORY_MINUTETIME_DATA);
        assert_eq!(pkg.req_header.zip_flag, 0x0c);
    }
}