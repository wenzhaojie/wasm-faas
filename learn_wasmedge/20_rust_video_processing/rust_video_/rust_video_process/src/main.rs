#[warn(unused_imports)]
extern crate serde;
extern crate serde_json;
extern crate image;

use serde::{Serialize, Deserialize};
use image::{GenericImageView};

// 定义图像类型
#[derive(Debug, Serialize, Deserialize)]
struct Image {
    width: u32,
    height: u32,
    pixels: Vec<Vec<(u8, u8, u8, u8)>>, // 修改为 Vec<Vec<(u8, u8, u8, u8)>>
}

fn resize(image_str: &str, width: u32, height: u32) -> String {
    let img: Image = serde_json::from_str(image_str).unwrap();
    let mut resized_pixels = Vec::new();

    for y in 0..height {
        let scaled_y = y * img.height / height;
        let mut row = Vec::new();
        for x in 0..width {
            let scaled_x = x * img.width / width;
            if let Some(px) = img.pixels.get(scaled_y as usize).and_then(|row| row.get(scaled_x as usize)) {
                row.push(*px);
            } else {
                // If out of bounds, push a default pixel value
                row.push((0, 0, 0, 255)); // Default black pixel
            }
        }
        resized_pixels.push(row);
    }

    serde_json::to_string(&Image {
        width,
        height,
        pixels: resized_pixels,
    }).unwrap()
}


// add_watermark 函数，接受图像字符串和水印字符串，返回添加水印后的图像字符串
fn add_watermark(image_str: &str, watermark_str: &str) -> String {
    let mut img: Image = serde_json::from_str(image_str).unwrap();
    let watermark_img: Image = serde_json::from_str(watermark_str).unwrap();
    let (wm_width, wm_height) = (watermark_img.width, watermark_img.height);
    let (img_width, img_height) = (img.width, img.height);

    // 打印相关信息
    println!("Image dimensions: {} x {}", img_width, img_height);
    println!("Watermark dimensions: {} x {}", wm_width, wm_height);

    let mut new_pixels = img.pixels.clone(); // 克隆原始图像的像素
    let offset_x = (img_width - wm_width) / 2;
    let offset_y = (img_height - wm_height) / 2;
    for y in 0..wm_height {
        for x in 0..wm_width {
            let (r, g, b, a) = watermark_img.pixels[y as usize][x as usize];
            if a > 0 {
                let dest_x = x + offset_x;
                let dest_y = y + offset_y;
                if dest_x < img_width && dest_y < img_height {
                    let px = &mut new_pixels[dest_y as usize][dest_x as usize];
                    px.0 = ((r as u32 + g as u32 + b as u32) / 3) as u8;
                    px.1 = ((r as u32 + g as u32 + b as u32) / 3) as u8;
                    px.2 = ((r as u32 + g as u32 + b as u32) / 3) as u8;
                }
            }
        }
    }
    img.pixels = new_pixels;
    serde_json::to_string(&img).unwrap()
}

// main 函数，程序入口
fn main() {
    // 示例使用
    // 从本地加载示例图像
    let original_image_path = "./original.png";
    let original_image = image::open(original_image_path).unwrap();
    let original_image_dimensions = original_image.dimensions();
    let original_image_str = serde_json::to_string(&Image {
        width: original_image_dimensions.0,
        height: original_image_dimensions.1,
        pixels: original_image.to_rgba8().pixels().map(|p| {
            let rgba = p.0;
            vec![(rgba[0], rgba[1], rgba[2], rgba[3])]
        }).collect(),
    }).unwrap();

    // 从本地加载示例水印图像
    let watermark_image_path = "./watermark.png";
    let watermark_image = image::open(watermark_image_path).unwrap();
    let watermark_image_dimensions = watermark_image.dimensions();
    let watermark_image_str = serde_json::to_string(&Image {
        width: watermark_image_dimensions.0,
        height: watermark_image_dimensions.1,
        pixels: watermark_image.to_rgba8().pixels().map(|p| {
            let rgba = p.0;
            vec![(rgba[0], rgba[1], rgba[2], rgba[3])]
        }).collect(),
    }).unwrap();

    let resized_image_str = resize(&original_image_str, 1000, 1000);
    println!("Resized Image: {}", resized_image_str);

    let watermarked_image_str = add_watermark(&resized_image_str, &watermark_image_str);
    println!("Watermarked Image: {}", watermarked_image_str);
}
