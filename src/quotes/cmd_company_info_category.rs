use std::io::{Cursor, Read};

use anyhow::{anyhow, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::internal::{sequence_id, utf8_to_gbk};
use crate::proto;
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};

#[derive(Debug, Default, Clone, Copy)]
pub struct CompanyInfoCategoryRequest {
    pub market: u16,
    pub code: [u8; 6],
    pub unknown: u32,
}

#[derive(Debug, Default, Clone)]
pub struct CompanyInfoCategory {
    pub name: String,
    pub filename: String,
    pub offset: u32,
    pub length: u32,
}

pub struct CompanyInfoCategoryPackage {
    req_header: StdRequestHeader,
    request: CompanyInfoCategoryRequest,
    resp_header: Option<StdResponseHeader>,
    reply: Vec<CompanyInfoCategory>,
}

impl CompanyInfoCategoryPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = proto::FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = proto::STD_MSG_COMPANY_CATEGORY;
        Ok(Self {
            req_header,
            request: CompanyInfoCategoryRequest::default(),
            resp_header: None,
            reply: Vec::new(),
        })
    }

    pub fn set_params(&mut self, request: CompanyInfoCategoryRequest) {
        self.request = request;
    }
}

impl Message for CompanyInfoCategoryPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 14u16; // 固定长度: 2 + 6 + 4 + 2(消息尾)
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        out.write_u16::<LittleEndian>(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        out.write_u32::<LittleEndian>(self.request.unknown)?;
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        if body.len() < 2 {
            return Err(anyhow!(
                "company info category response too short: {}",
                body.len()
            ));
        }
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);
        let count = cursor.read_u16::<LittleEndian>()? as usize;
        let expected_len = 2 + count * 152;
        if body.len() < expected_len {
            return Err(anyhow!(
                "company info category truncated: have={}, expect>= {}",
                body.len(),
                expected_len
            ));
        }
        self.reply.clear();
        self.reply.reserve(count);
        for _ in 0..count {
            let mut name_buf = [0u8; 64];
            cursor.read_exact(&mut name_buf)?;
            let mut filename_buf = [0u8; 80];
            cursor.read_exact(&mut filename_buf)?;
            let offset = cursor.read_u32::<LittleEndian>()?;
            let length = cursor.read_u32::<LittleEndian>()?;
            let name = utf8_to_gbk(&name_buf);
            let filename = utf8_to_gbk(&filename_buf);
            self.reply.push(CompanyInfoCategory {
                name,
                filename,
                offset,
                length,
            });
        }
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_company_info_category() -> Result<CompanyInfoCategoryPackage> {
    CompanyInfoCategoryPackage::new()
}
