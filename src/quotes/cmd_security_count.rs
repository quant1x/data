use anyhow::{anyhow, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::io::Cursor;

use crate::internal::{hex_string_to_bytes, sequence_id};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default, Clone, Copy)]
pub struct SecurityCountRequest {
    pub market: u16,
}

#[derive(Debug, Default, Clone, Copy)]
pub struct SecurityCountReply {
    pub count: u16,
}

pub struct SecurityCountPackage {
    req_header: StdRequestHeader,
    request: SecurityCountRequest,
    resp_header: Option<StdResponseHeader>,
    reply: SecurityCountReply,
    content: Vec<u8>,
}

impl SecurityCountPackage {
    pub fn new() -> Result<Self> {
        let content = hex_string_to_bytes("75c73301")?;
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_SECURITY_COUNT;
        Ok(Self {
            req_header,
            request: SecurityCountRequest { market: 0 },
            resp_header: None,
            reply: SecurityCountReply::default(),
            content,
        })
    }

    pub fn set_params(&mut self, request: SecurityCountRequest) {
        self.request = request;
    }
}

impl Message for SecurityCountPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len =
            2 + std::mem::size_of::<SecurityCountRequest>() as u16 + self.content.len() as u16;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.content);
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        if body.len() < 2 {
            return Err(anyhow!("security count response too short: {}", body.len()));
        }
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);
        self.reply.count = cursor.read_u16::<LittleEndian>()?;
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_security_count() -> Result<SecurityCountPackage> {
    SecurityCountPackage::new()
}
