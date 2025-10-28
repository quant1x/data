/// 判断 `f64` 是否为 NaN 或无穷。
pub fn is_nan_or_inf(value: f64) -> bool {
    value.is_nan() || value.is_infinite()
}

/// 泛型数值转 `f64`。
pub fn number_to_f64<T>(value: T) -> f64
where
    T: Into<f64>,
{
    value.into()
}

/// 按照 Go 版本 `IntToFloat64` 的算法将整数解码为 `f64`。
pub fn int_to_f64<T>(integer: T) -> f64
where
    T: Into<i64>,
{
    let ivol = integer.into() as i64;
    let log_point = (ivol >> 24) as i32;
    let hleax = ((ivol >> 16) & 0xff) as i32;
    let lheax = ((ivol >> 8) & 0xff) as i32;
    let lleax = (ivol & 0xff) as i32;

    let dw_ecx = log_point * 2 - 0x7f;
    let dw_edx = log_point * 2 - 0x86;
    let dw_esi = log_point * 2 - 0x8e;
    let dw_eax = log_point * 2 - 0x96;

    let tmp_eax = if dw_ecx < 0 { -dw_ecx } else { dw_ecx };
    let mut dbl_xmm6 = 2f64.powf(tmp_eax as f64);
    if dw_ecx < 0 {
        dbl_xmm6 = 1.0 / dbl_xmm6;
    }

    let dbl_xmm4 = if hleax > 0x80 {
        let tmpdbl_xmm3 = 2f64.powf((dw_edx + 1) as f64);
        let mut dbl_xmm0 = 2f64.powf(dw_edx as f64) * 128.0;
        dbl_xmm0 += (hleax & 0x7f) as f64 * tmpdbl_xmm3;
        dbl_xmm0
    } else {
        if dw_edx >= 0 {
            2f64.powf(dw_edx as f64) * hleax as f64
        } else {
            (1.0 / 2f64.powf(dw_edx as f64)) * hleax as f64
        }
    };

    let mut dbl_xmm3 = 2f64.powf(dw_esi as f64) * lheax as f64;
    let mut dbl_xmm1 = 2f64.powf(dw_eax as f64) * lleax as f64;
    if hleax & 0x80 != 0 {
        dbl_xmm3 *= 2.0;
        dbl_xmm1 *= 2.0;
    }

    dbl_xmm6 + dbl_xmm4 + dbl_xmm3 + dbl_xmm1
}
