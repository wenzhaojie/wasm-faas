use image::{DynamicImage, GenericImageView, ImageBuffer, Rgba, ImageOutputFormat, open};
use serde::{Serialize, Deserialize};
use std::fs::File;



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

fn main() {
    // Read the image from a file called "input.jpg"
    let img = open("input.jpg").expect("Failed to open input.jpg");
    let original_image = serialize_image(&img);

    let processed_image = apply_lut(&original_image);
    let dynamic_image = deserialize_image(&processed_image);

    // Save the processed image to a file
    let mut output_file = File::create("output.jpg").unwrap();
    dynamic_image.write_to(&mut output_file, ImageOutputFormat::Jpeg(80)).unwrap();

    println!("Processed image saved to 'output.jpg'.");
}
