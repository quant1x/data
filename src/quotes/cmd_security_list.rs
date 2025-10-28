use std::io::{Cursor, Read};

use anyhow::{anyhow, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::internal::{int_to_f64, sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default, Clone, Copy)]
pub struct SecurityListRequest {
    pub market: u16,
    pub start: u16,
}

#[derive(Debug, Default, Clone)]
pub struct Security {
    pub code: String,
    pub vol_unit: u16,
    pub reversed1: [u8; 4],
    pub decimal_point: i8,
    pub name: String,
    pub pre_close: f64,
    pub reversed2: [u8; 4],
}

#[derive(Debug, Default, Clone)]
pub struct SecurityListReply {
    pub count: u16,
    pub list: Vec<Security>,
}

pub struct SecurityListPackage {
    req_header: StdRequestHeader,
    request: SecurityListRequest,
    resp_header: Option<StdResponseHeader>,
    reply: SecurityListReply,
}

impl SecurityListPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_SECURITY_LIST;
        Ok(Self {
            req_header,
            request: SecurityListRequest::default(),
            resp_header: None,
            reply: SecurityListReply::default(),
        })
    }

    pub fn set_params(&mut self, request: SecurityListRequest) {
        self.request = request;
    }
}

impl Message for SecurityListPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 2 + std::mem::size_of::<SecurityListRequest>() as u16;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.write_u16::<LittleEndian>(self.request.start)?;
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        if body.len() < 2 {
            return Err(anyhow!("security list response too short: {}", body.len()));
        }
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;
        self.reply.list.clear();

        const ENTRY_SIZE: usize = 6 + 2 + 8 + 4 + 1 + 4 + 4;
        for _ in 0..self.reply.count {
            let remaining = body.len().saturating_sub(cursor.position() as usize);
            if remaining < ENTRY_SIZE {
                return Err(anyhow!(
                    "security list entry truncated: remaining={remaining}"
                ));
            }

            let mut code_bytes = [0u8; 6];
            cursor.read_exact(&mut code_bytes)?;
            let code = String::from_utf8_lossy(&code_bytes)
                .trim_end_matches('\u{0}')
                .to_string();

            let vol_unit = cursor.read_u16::<LittleEndian>()?;

            let mut name_bytes = [0u8; 8];
            cursor.read_exact(&mut name_bytes)?;
            let name = utf8_to_gbk(&name_bytes);

            let mut reversed1 = [0u8; 4];
            cursor.read_exact(&mut reversed1)?;

            let decimal_point = cursor.read_i8()?;

            let raw_pre_close = cursor.read_u32::<LittleEndian>()?;
            let pre_close = int_to_f64(raw_pre_close);

            let mut reversed2 = [0u8; 4];
            cursor.read_exact(&mut reversed2)?;

            self.reply.list.push(Security {
                code,
                vol_unit,
                reversed1,
                decimal_point,
                name,
                pre_close,
                reversed2,
            });
        }
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_security_list() -> Result<SecurityListPackage> {
    SecurityListPackage::new()
}
