use std::slice;
use std::time::Instant;
use std::convert::TryInto; // 添加导入 TryInto trait

#[no_mangle]
pub extern "C" fn bandwidth(data: *const i32, length: i32) -> f64 {
    let start_time = Instant::now();

    // Convert length from i32 to usize
    let usize_length: usize = length.try_into().unwrap(); // 尝试转换为 usize

    // Read data from memory
    let _data_slice = unsafe { slice::from_raw_parts(data, usize_length) };

    // Process the data (optional)

    let elapsed_time = start_time.elapsed();
    let elapsed_secs = elapsed_time.as_secs_f64();

    elapsed_secs
}
