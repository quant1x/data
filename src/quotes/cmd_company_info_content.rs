use std::io::{Cursor, Read};

use anyhow::{anyhow, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::internal::{sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Clone, Copy)]
pub struct CompanyInfoContentRequest {
    pub market: u16,
    pub code: [u8; 6],
    pub unknown1: u16,
    pub filename: [u8; 80],
    pub offset: u32,
    pub length: u32,
    pub unknown2: u32,
}

impl Default for CompanyInfoContentRequest {
    fn default() -> Self {
        Self {
            market: 0,
            code: [0u8; 6],
            unknown1: 0,
            filename: [0u8; 80],
            offset: 0,
            length: 0,
            unknown2: 0,
        }
    }
}

#[derive(Debug, Default, Clone)]
pub struct CompanyInfoContent {
    pub market: u16,
    pub code: String,
    pub name: Option<String>,
    pub length: u32,
    pub content: String,
}

pub struct CompanyInfoContentPackage {
    req_header: StdRequestHeader,
    request: CompanyInfoContentRequest,
    resp_header: Option<StdResponseHeader>,
    reply: CompanyInfoContent,
}

impl CompanyInfoContentPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_COMPANY_CONTENT;
        Ok(Self {
            req_header,
            request: CompanyInfoContentRequest::default(),
            resp_header: None,
            reply: CompanyInfoContent::default(),
        })
    }

    pub fn set_params(&mut self, request: CompanyInfoContentRequest) {
        self.request = request;
    }
}

impl Message for CompanyInfoContentPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let request_len = 2 + 6 + 2 + 80 + 4 + 4 + 4;
        let payload_len = 2 + request_len as u16;
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u16::<LittleEndian>(self.request.unknown1)?;
        out.extend_from_slice(&self.request.filename);
        out.write_u32::<LittleEndian>(self.request.offset)?;
        out.write_u32::<LittleEndian>(self.request.length)?;
        out.write_u32::<LittleEndian>(self.request.unknown2)?;
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        if body.len() < 12 {
            return Err(anyhow!(
                "company info content response too short: {}",
                body.len()
            ));
        }
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);
        let market = cursor.read_u16::<LittleEndian>()?;
        let mut code_bytes = [0u8; 6];
        cursor.read_exact(&mut code_bytes)?;
        let mut _unknown1 = [0u8; 2];
        cursor.read_exact(&mut _unknown1)?;
        let length = cursor.read_u16::<LittleEndian>()?;
        let mut data = vec![0u8; length as usize];
        cursor.read_exact(&mut data)?;

        self.reply.market = market;
        self.reply.code = String::from_utf8_lossy(&code_bytes)
            .trim_end_matches('\u{0}')
            .to_string();
        self.reply.name = None;
        self.reply.length = length as u32;
        self.reply.content = utf8_to_gbk(&data);
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_company_info_content() -> Result<CompanyInfoContentPackage> {
    CompanyInfoContentPackage::new()
}
