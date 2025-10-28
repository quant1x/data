use std::any::Any;
use std::io::{Cursor, Read};

use anyhow::{anyhow, Context, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use flate2::read::ZlibDecoder;
use log::{debug, warn};

use crate::quotes::base_client::TcpClient;
use crate::quotes::base_consts::{QuotesError, MESSAGE_HEADER_BYTES, MESSAGE_MAX_BYTES};

#[derive(Debug, Default, Clone, Copy)]
pub struct StdRequestHeader {
    pub zip_flag: u8,
    pub seq_id: u32,
    pub packet_type: u8,
    pub pkg_len1: u16,
    pub pkg_len2: u16,
    pub method: u16,
}

impl StdRequestHeader {
    pub fn to_bytes(&self) -> Result<Vec<u8>> {
        let mut buf = Vec::with_capacity(12);
        buf.write_u8(self.zip_flag)?;
        buf.write_u32::<LittleEndian>(self.seq_id)?;
        buf.write_u8(self.packet_type)?;
        buf.write_u16::<LittleEndian>(self.pkg_len1)?;
        buf.write_u16::<LittleEndian>(self.pkg_len2)?;
        buf.write_u16::<LittleEndian>(self.method)?;
        Ok(buf)
    }
}

#[derive(Debug, Default, Clone, Copy)]
pub struct StdResponseHeader {
    pub i1: u32,
    pub zip_flag: u8,
    pub seq_id: u32,
    pub i3: u8,
    pub method: u16,
    pub zip_size: u16,
    pub unzip_size: u16,
}

impl StdResponseHeader {
    fn from_bytes(bytes: &[u8]) -> Result<Self> {
        if bytes.len() < MESSAGE_HEADER_BYTES {
            return Err(anyhow!("response header too short: {}", bytes.len()));
        }
        let mut cursor = Cursor::new(bytes);
        Ok(Self {
            i1: cursor.read_u32::<LittleEndian>()?,
            zip_flag: cursor.read_u8()?,
            seq_id: cursor.read_u32::<LittleEndian>()?,
            i3: cursor.read_u8()?,
            method: cursor.read_u16::<LittleEndian>()?,
            zip_size: cursor.read_u16::<LittleEndian>()?,
            unzip_size: cursor.read_u16::<LittleEndian>()?,
        })
    }
}

pub trait Message {
    fn serialize(&mut self) -> Result<Vec<u8>>;
    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()>;
    fn reply(&self) -> &(dyn Any + Send + Sync);
}

pub fn process(client: &TcpClient, msg: &mut dyn Message) -> Result<()> {
    let send_data = msg.serialize()?;
    if send_data.is_empty() {
        return Err(anyhow!("message serialized to empty payload"));
    }
    debug!("send data len={}", send_data.len());
    client.write_all(&send_data)?;

    let mut header_bytes = [0u8; MESSAGE_HEADER_BYTES];
    client.read_exact(&mut header_bytes)?;
    debug!("response header bytes={:x?}", header_bytes);

    let header = StdResponseHeader::from_bytes(&header_bytes)?;
    debug!("response header={:?}", header);

    let zip_size = header.zip_size as usize;
    let unzip_size = header.unzip_size as usize;
    if zip_size > MESSAGE_MAX_BYTES {
        warn!(
            "response size {} exceeds max {}",
            zip_size, MESSAGE_MAX_BYTES
        );
        return Err(QuotesError::BadData.into());
    }

    let mut msg_data = vec![0u8; zip_size];
    client.read_exact(&mut msg_data)?;

    let body = if zip_size != unzip_size {
        let mut decoder = ZlibDecoder::new(&msg_data[..]);
        let mut out = Vec::with_capacity(unzip_size.max(zip_size));
        decoder
            .read_to_end(&mut out)
            .context("decompress response body")?;
        out
    } else {
        msg_data
    };

    msg.unserialize(&header, &body)
}
