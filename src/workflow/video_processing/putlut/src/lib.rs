// 这是一个handler，用于使用wasm环境沙箱来运行FaaS函数
use serde_json;
use redis::Commands;
use std::time::Instant;
use std::collections::HashMap;
use wasmedge_bindgen_macro::*;

use image::{DynamicImage, GenericImageView, ImageBuffer, Rgba, open};
use serde::{Serialize, Deserialize};


#[derive(Serialize, Deserialize, Clone)]
struct SerializableImage {
    width: u32,
    height: u32,
    pixels: Vec<u8>,
}

fn get_url() -> String {
    if let Ok(url) = std::env::var("REDIS_URL") {
        url
    } else {
        "redis://127.0.0.1/".into()
    }
}


fn deserialize_image(serializable_image: &SerializableImage) -> DynamicImage {
    let SerializableImage { width, height, pixels } = serializable_image;
    let image_buffer = ImageBuffer::<Rgba<u8>, _>::from_raw(*width, *height, pixels.clone()).unwrap();
    DynamicImage::ImageRgba8(image_buffer)
}


fn serialize_image(dynamic_image: &DynamicImage) -> SerializableImage {
    let (width, height) = dynamic_image.dimensions();
    let raw_pixels = dynamic_image.to_rgba8().into_raw();
    SerializableImage {
        width,
        height,
        pixels: raw_pixels,
    }
}

#[wasmedge_bindgen]
pub fn helloworld(input_str: String) -> String {
    return format!("Hello, {}!", input_str);
}



#[wasmedge_bindgen]
pub fn put_input_img_into_redis(input_img_path: String, input_obj_key: String) -> String {
    let img = open(input_img_path).unwrap();
    let serializable_image = serialize_image(&img);
    let serializable_image_str = serde_json::to_string(&serializable_image).unwrap();
    let client = redis::Client::open(&*get_url()).unwrap();
    let mut con = client.get_connection().unwrap();
    let _ : () = con.set(input_obj_key, &serializable_image_str).unwrap();
    // 返回一个字符串，表示成功
    "Success".to_string()
}


fn apply_lut(serializable_image: &SerializableImage) -> SerializableImage {
    let mut image = deserialize_image(serializable_image);

    if let DynamicImage::ImageRgba8(ref mut img) = image {
        let (width, height) = img.dimensions();
        let mut imgbuf: ImageBuffer<Rgba<u8>, Vec<u8>> = ImageBuffer::new(width, height);

        for (x, y, pixel) in img.enumerate_pixels() {
            let mut data = *pixel; // Destructure to get a copy of the color
            // Update each channel directly, increment by 50 or up to the maximum value
            data.0[0] = data.0[0].saturating_add(50);
            data.0[1] = data.0[1].saturating_add(50);
            data.0[2] = data.0[2].saturating_add(50);
            imgbuf.put_pixel(x, y, data); // Use the modified data directly
        }

        let processed_image = DynamicImage::ImageRgba8(imgbuf);
        serialize_image(&processed_image)
    } else {
        panic!("Unsupported image format for LUT application");
    }
}


#[wasmedge_bindgen]
pub fn handler(input_obj_key: String, output_obj_key: String) -> String {
    let start_conn_redis_t = Instant::now();
    let client = redis::Client::open(&*get_url()).unwrap();
    let mut con = client.get_connection().unwrap();
    let conn_redis_t = start_conn_redis_t.elapsed();

    // 计时
    let start_get_input_t = Instant::now();
    let input_obj_serde_str: String = con.get(&input_obj_key).unwrap(); // Get the serialized string from Redis
    let get_input_t = start_get_input_t.elapsed();

    // 反序列化输入字符串为输入对象
    let input_obj: SerializableImage = serde_json::from_str(&input_obj_serde_str).unwrap(); // Deserialize into an owned SerializableImage
    let start_compute_t = Instant::now();
    let output_obj = apply_lut(&input_obj); // Passed as reference here if needed, otherwise pass owned
    let compute_t = start_compute_t.elapsed();

    // Serialize output object
    let output_obj_serde_str = serde_json::to_string(&output_obj).unwrap();

    // 计时
    let start_set_output_t = Instant::now();
    let _ : () = con.set(&output_obj_key, &output_obj_serde_str).unwrap();
    let set_output_t = start_set_output_t.elapsed();

    // 返回统计信息
    let mut stats_dict = HashMap::new();
    stats_dict.insert("conn_redis_t", conn_redis_t.as_millis());
    stats_dict.insert("get_input_t", get_input_t.as_millis());
    stats_dict.insert("compute_t", compute_t.as_millis());
    stats_dict.insert("set_output_t", set_output_t.as_millis());

    serde_json::to_string(&stats_dict).unwrap()
}