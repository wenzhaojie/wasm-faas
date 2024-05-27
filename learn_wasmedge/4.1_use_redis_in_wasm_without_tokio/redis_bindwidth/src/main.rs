use redis::Commands;
use anyhow::Result;
use std::time::{Instant};


const DATA_SIZE_MB: f64 = 300.0; // Data size in MB

fn get_url() -> String {
    if let Ok(url) = std::env::var("REDIS_URL") {
        url
    } else {
        "redis://127.0.0.1/".into()
    }
}

fn main() -> Result<(), E> {
    // Connect to Redis
    let client = redis::Client::open(&*get_url())?;
    let mut con = client.get_connection()?;

    // Generate data

    let data: String = (0..(DATA_SIZE_MB * 1024.0 * 1024.0 / 8.0) as usize).map(|_| "01234567").collect();

    // Start timer
    let start_time = Instant::now();

    // Set key-value pair in Redis
    con.set("large_data", &data)?;

    // Calculate elapsed time
    let elapsed = start_time.elapsed();
    println!("Time taken to set {:.2} MB data in Redis: {:?}", DATA_SIZE_MB, elapsed);

    // Calculate bandwidth
    let elapsed_seconds = elapsed.as_secs() as f64 + elapsed.subsec_nanos() as f64 / 1_000_000_000.0; // Elapsed time in seconds
    let bandwidth = DATA_SIZE_MB / elapsed_seconds;
    println!("Bandwidth: {:.2} MB/s", bandwidth);

    Ok(())
}

