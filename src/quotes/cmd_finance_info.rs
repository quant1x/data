use crate::internal::sequence_id;
use crate::proto::{FLAG_NOT_ZIPPED, STD_MSG_FINANCE_INFO};
use crate::quotes::message::{Message, StdRequestHeader, StdResponseHeader};
use anyhow::{anyhow, Result};
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::io::{Cursor, Read};

#[derive(Debug, Clone)]
pub struct FinanceInfoRequest {
    pub count: u16,
    pub market: u8,
    pub code: [u8; 6],
}

#[derive(Debug, Clone)]
pub struct FinanceInfo {
    pub code: String,
    pub liu_tong_gu_ben: f64,
    pub province: u16,
    pub industry: u16,
    pub updated_date: u32,
    pub ipo_date: u32,
    pub zong_gu_ben: f64,
    pub guo_jia_gu: f64,
    pub fa_qi_ren_fa_ren_gu: f64,
    pub fa_ren_gu: f64,
    pub b_gu: f64,
    pub h_gu: f64,
    pub zhi_gong_gu: f64,
    pub zong_zi_chan: f64,
    pub liu_dong_zi_chan: f64,
    pub gu_ding_zi_chan: f64,
    pub wu_xing_zi_chan: f64,
    pub gu_dong_ren_shu: f64,
    pub liu_dong_fu_zhai: f64,
    pub chang_qi_fu_zhai: f64,
    pub zi_ben_gong_ji_jin: f64,
    pub jing_zi_chan: f64,
    pub zhu_ying_shou_ru: f64,
    pub zhu_ying_li_run: f64,
    pub ying_shou_zhang_kuan: f64,
    pub ying_ye_li_run: f64,
    pub tou_zi_shou_yu: f64,
    pub jing_ying_xian_jin_liu: f64,
    pub zong_xian_jin_liu: f64,
    pub cun_huo: f64,
    pub li_run_zong_he: f64,
    pub shui_hou_li_run: f64,
    pub jing_li_run: f64,
    pub wei_fen_li_run: f64,
    pub mei_gu_jing_zi_chan: f64,
    pub bao_liu_2: f64,
}

impl Default for FinanceInfo {
    fn default() -> Self {
        FinanceInfo {
            code: String::new(),
            liu_tong_gu_ben: 0.0,
            province: 0,
            industry: 0,
            updated_date: 0,
            ipo_date: 0,
            zong_gu_ben: 0.0,
            guo_jia_gu: 0.0,
            fa_qi_ren_fa_ren_gu: 0.0,
            fa_ren_gu: 0.0,
            b_gu: 0.0,
            h_gu: 0.0,
            zhi_gong_gu: 0.0,
            zong_zi_chan: 0.0,
            liu_dong_zi_chan: 0.0,
            gu_ding_zi_chan: 0.0,
            wu_xing_zi_chan: 0.0,
            gu_dong_ren_shu: 0.0,
            liu_dong_fu_zhai: 0.0,
            chang_qi_fu_zhai: 0.0,
            zi_ben_gong_ji_jin: 0.0,
            jing_zi_chan: 0.0,
            zhu_ying_shou_ru: 0.0,
            zhu_ying_li_run: 0.0,
            ying_shou_zhang_kuan: 0.0,
            ying_ye_li_run: 0.0,
            tou_zi_shou_yu: 0.0,
            jing_ying_xian_jin_liu: 0.0,
            zong_xian_jin_liu: 0.0,
            cun_huo: 0.0,
            li_run_zong_he: 0.0,
            shui_hou_li_run: 0.0,
            jing_li_run: 0.0,
            wei_fen_li_run: 0.0,
            mei_gu_jing_zi_chan: 0.0,
            bao_liu_2: 0.0,
        }
    }
}

impl FinanceInfo {
    pub fn is_delisting(&self) -> bool {
        self.ipo_date == 0 && self.zong_gu_ben == 0.0 && self.liu_tong_gu_ben == 0.0
    }
}

pub struct FinanceInfoPackage {
    req_header: StdRequestHeader,
    request: FinanceInfoRequest,
    resp_header: Option<StdResponseHeader>,
    reply: FinanceInfo,
}

impl FinanceInfoPackage {
    pub fn new() -> Result<Self> {
        let mut req_header = StdRequestHeader::default();
        req_header.zip_flag = FLAG_NOT_ZIPPED;
        req_header.packet_type = 0x01;
        req_header.method = STD_MSG_FINANCE_INFO;
        Ok(Self {
            req_header,
            request: FinanceInfoRequest {
                count: 1,
                market: 0,
                code: [0; 6],
            },
            resp_header: None,
            reply: FinanceInfo::default(),
        })
    }

    pub fn set_params(&mut self, request: FinanceInfoRequest) {
        self.request = request;
    }
}

impl Message for FinanceInfoPackage {
    fn serialize(&mut self) -> Result<Vec<u8>> {
        self.req_header.seq_id = sequence_id();
        let payload_len = 2 + 1 + 6; // count(u16) + market(u8) + code([u8;6])
        self.req_header.pkg_len1 = payload_len;
        self.req_header.pkg_len2 = payload_len;

        let mut out = self.req_header.to_bytes()?;
        out.write_u16::<LittleEndian>(self.request.count)?;
        out.write_u8(self.request.market)?;
        out.extend_from_slice(&self.request.code);
        Ok(out)
    }

    fn unserialize(&mut self, header: &StdResponseHeader, body: &[u8]) -> Result<()> {
        if body.len() < 9 { // market(u8) + code([u8;6]) + at least some data
            return Err(anyhow!("finance info response too short: {}", body.len()));
        }
        self.resp_header = Some(*header);
        let mut cursor = Cursor::new(body);

        let _market = cursor.read_u8()?;
        let mut code = [0u8; 6];
        cursor.read_exact(&mut code)?;

        let mut reply = FinanceInfo::default();
        reply.code = String::from_utf8_lossy(&code).trim_end_matches('\0').to_string();

        reply.liu_tong_gu_ben = cursor.read_f32::<LittleEndian>()? as f64;
        reply.province = cursor.read_u16::<LittleEndian>()?;
        reply.industry = cursor.read_u16::<LittleEndian>()?;
        reply.updated_date = cursor.read_u32::<LittleEndian>()?;
        reply.ipo_date = cursor.read_u32::<LittleEndian>()?;
        reply.zong_gu_ben = cursor.read_f32::<LittleEndian>()? as f64;
        reply.guo_jia_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.fa_qi_ren_fa_ren_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.fa_ren_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.b_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.h_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zhi_gong_gu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zong_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.liu_dong_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.gu_ding_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.wu_xing_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.gu_dong_ren_shu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.liu_dong_fu_zhai = cursor.read_f32::<LittleEndian>()? as f64;
        reply.chang_qi_fu_zhai = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zi_ben_gong_ji_jin = cursor.read_f32::<LittleEndian>()? as f64;
        reply.jing_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zhu_ying_shou_ru = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zhu_ying_li_run = cursor.read_f32::<LittleEndian>()? as f64;
        reply.ying_shou_zhang_kuan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.ying_ye_li_run = cursor.read_f32::<LittleEndian>()? as f64;
        reply.tou_zi_shou_yu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.jing_ying_xian_jin_liu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.zong_xian_jin_liu = cursor.read_f32::<LittleEndian>()? as f64;
        reply.cun_huo = cursor.read_f32::<LittleEndian>()? as f64;
        reply.li_run_zong_he = cursor.read_f32::<LittleEndian>()? as f64;
        reply.shui_hou_li_run = cursor.read_f32::<LittleEndian>()? as f64;
        reply.jing_li_run = cursor.read_f32::<LittleEndian>()? as f64;
        reply.wei_fen_li_run = cursor.read_f32::<LittleEndian>()? as f64;
        reply.mei_gu_jing_zi_chan = cursor.read_f32::<LittleEndian>()? as f64;
        reply.bao_liu_2 = cursor.read_f32::<LittleEndian>()? as f64;

        self.reply = reply;
        Ok(())
    }

    fn reply(&self) -> &(dyn std::any::Any + Send + Sync) {
        &self.reply
    }
}

pub fn new_finance_info() -> Result<FinanceInfoPackage> {
    FinanceInfoPackage::new()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_finance_info_default() {
        let info = FinanceInfo::default();
        assert_eq!(info.code, "");
        assert_eq!(info.liu_tong_gu_ben, 0.0);
        assert_eq!(info.ipo_date, 0);
    }

    #[test]
    fn test_finance_info_is_delisting() {
        let mut info = FinanceInfo::default();
        assert!(info.is_delisting());

        info.ipo_date = 1;
        assert!(!info.is_delisting());

        info.ipo_date = 0;
        info.zong_gu_ben = 1.0;
        assert!(!info.is_delisting());
    }
}