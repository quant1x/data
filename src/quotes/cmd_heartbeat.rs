use anyhow::Result;

use crate::internal::{hex_string_to_bytes, sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default, Clone)]
pub struct HeartBeatReply {
    pub info: String,
}

impl HeartBeatReply {
    pub fn is_empty(&self) -> bool {
        self.info.is_empty()
    }
}

pub struct HeartBeatPackage {
    req_header: StdRequestHeader,
    resp_header: Option<StdResponseHeader>,
    reply: HeartBeatReply,
    content: Vec<u8>,
}

impl HeartBeatPackage {
    pub fn new() -> Result<Self> {
        let content = hex_string_to_bytes("02020002000400")?;
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x02;
        req_header.method = proto::STD_MSG_HEARTBEAT;
        Ok(Self {
            req_header,
            resp_header: None,
            reply: HeartBeatReply::default(),
            content,
        })
    }
}

impl Message for HeartBeatPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = self.content.len() as u16 + 2;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;
        let mut out = self.req_header.to_bytes()?;
        out.extend_from_slice(&self.content);
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        self.resp_header = Some(*header);
        self.reply.info = utf8_to_gbk(body);
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_heartbeat() -> Result<HeartBeatPackage> {
    HeartBeatPackage::new()
}
