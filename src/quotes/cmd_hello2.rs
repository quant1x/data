use anyhow::Result;

use crate::internal::{hex_string_to_bytes, sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default, Clone)]
pub struct Hello2Reply {
    pub info: String,
    pub server_time: Option<String>,
}

pub struct Hello2Package {
    req_header: StdRequestHeader,
    resp_header: Option<StdResponseHeader>,
    reply: Hello2Reply,
    content: Vec<u8>,
}

impl Hello2Package {
    pub fn new() -> Result<Self> {
        let content =
            hex_string_to_bytes("d5d0c9ccd6a4a8af0000008fc22540130000d500c9ccbdf0d7ea00000002")?;
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_LOGIN2;
        Ok(Self {
            req_header,
            resp_header: None,
            reply: Hello2Reply::default(),
            content,
        })
    }
}

impl Message for Hello2Package {
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
        if body.len() >= 58 {
            self.reply.info = utf8_to_gbk(&body[58..]);
        } else {
            self.reply.info = utf8_to_gbk(body);
        }
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_hello2() -> Result<Hello2Package> {
    Hello2Package::new()
}
