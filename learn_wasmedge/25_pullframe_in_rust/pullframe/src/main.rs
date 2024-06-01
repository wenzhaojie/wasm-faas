extern crate opencv;
extern crate image;
extern crate serde;

use opencv::{
    prelude::*,
    videoio,
    core,
};
use serde::{Serialize, Deserialize};
use image::{DynamicImage, ImageBuffer, Rgba, GenericImageView};

#[derive(Serialize, Deserialize, Clone)]
struct SerializableImage {
    width: u32,
    height: u32,
    pixels: Vec<u8>,
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

fn extract_frames_from_video(file_path: &str) -> opencv::Result<Vec<SerializableImage>> {
    let mut video_capture = videoio::VideoCapture::from_file(file_path, videoio::CAP_ANY)?;
    let mut frames = Vec::new();

    let mut frame = core::Mat::default();
    while video_capture.read(&mut frame)? {
        if frame.size()?.width > 0 {
            let mut buf = core::Vector::<u8>::new();
            let params = core::Vector::<i32>::new();
            let success = opencv::imgcodecs::imencode(".png", &frame, &mut buf, &params)?;
            if !success {
                return Err(opencv::Error::new(opencv::core::StsError, "Failed to encode image".to_string()));
            }
            let image = image::load_from_memory(buf.as_slice())
                .map_err(|e| opencv::Error::new(opencv::core::StsError, format!("Image load error: {}", e)))?
                .to_rgba8();
            let serializable_image = serialize_image(&DynamicImage::ImageRgba8(image));
            frames.push(serializable_image);
        }
    }

    Ok(frames)
}

fn main() {
    match extract_frames_from_video("path/to/your/video.mp4") {
        Ok(frames) => println!("Extracted {} frames.", frames.len()),
        Err(e) => println!("Failed to extract frames: {}", e),
    }
}
