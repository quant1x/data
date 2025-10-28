use anyhow::Result;

use crate::internal::{hex_string_to_bytes, sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default)]
pub struct Hello1Request;

#[derive(Debug, Default, Clone)]
pub struct Hello1Reply {
    pub info: String,
    pub server_time: Option<String>,
}

pub struct Hello1Package {
    req_header: StdRequestHeader,
    _request: Hello1Request,
    resp_header: Option<StdResponseHeader>,
    reply: Hello1Reply,
    content: Vec<u8>,
}

impl Hello1Package {
    pub fn new() -> Result<Self> {
        let content = hex_string_to_bytes("01")?;
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_LOGIN1;
        Ok(Self {
            req_header,
            _request: Hello1Request::default(),
            resp_header: None,
            reply: Hello1Reply::default(),
            content,
        })
    }
}

impl Message for Hello1Package {
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
        if body.len() < 68 {
            self.reply.info = utf8_to_gbk(body);
        } else {
            let info_bytes = &body[68..];
            self.reply.info = utf8_to_gbk(info_bytes);
        }
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_hello1() -> Result<Hello1Package> {
    Hello1Package::new()
}
