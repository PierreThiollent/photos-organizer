use chrono::{DateTime, Local};
use std::env;
use std::io::{Error, Result};
use std::time::Instant;

fn main() -> Result<()> {
    let mut moved_count = 0usize;
    let mut error_count = 0usize;
    let mut folder_count = 0usize;

    let executable_path = env::current_exe()?;
    let working_dir = executable_path
        .parent()
        .ok_or_else(|| Error::other("Failed to get parent directory"))?
        .to_path_buf();

    println!("🏞️ Photos/videos organizer");
    println!("Processing files...\n");

    let started_at = Instant::now();

    for entry in working_dir.read_dir()? {
        let entry = match entry {
            Ok(entry) => entry,
            Err(error) => {
                error_count += 1;
                eprintln!("Failed to read entry: {error}");
                continue;
            }
        };

        let path = entry.path();
        let file_name = entry.file_name();

        if path == executable_path
            || path.extension().is_none()
            || file_name.as_encoded_bytes().starts_with(b".")
        {
            continue;
        }

        let metadata = match entry.metadata() {
            Ok(metadata) => metadata,
            Err(error) => {
                error_count += 1;
                eprintln!("Failed to read metadata for {}: {error}", path.display());
                continue;
            }
        };

        if !metadata.is_file() {
            continue;
        }

        let modification_time = match metadata.modified() {
            Ok(time) => time,
            Err(error) => {
                error_count += 1;
                eprintln!(
                    "Failed to get modification date for {}: {error}",
                    path.display()
                );
                continue;
            }
        };

        let modification_date = DateTime::<Local>::from(modification_time)
            .format("%Y-%m-%d")
            .to_string();

        let new_dir = working_dir.join(&modification_date);
        if !new_dir.exists() {
            if let Err(error) = std::fs::create_dir(&new_dir) {
                error_count += 1;
                eprintln!("Failed to create directory {}: {error}", new_dir.display());
                continue;
            }

            folder_count += 1;
        }

        let new_path = new_dir.join(file_name);
        if let Err(error) = std::fs::rename(&path, &new_path) {
            error_count += 1;
            eprintln!(
                "Failed to move file {} to {}: {error}",
                path.display(),
                new_path.display()
            );
            continue;
        }

        moved_count += 1;
    }

    let elapsed = started_at.elapsed();
    let elapsed_seconds = elapsed.as_secs_f64();
    let files_per_second = if elapsed_seconds > 0.0 {
        moved_count as f64 / elapsed_seconds
    } else {
        0.0
    };
    let status = if error_count == 0 {
        "Success"
    } else {
        "Completed with errors"
    };
    let duration = format!("{elapsed_seconds:.3} s");
    let throughput = format!("{files_per_second:.0} files/s");

    println!("┌──────────────────────────────────────────┐");
    println!("│ {:<40} │", "Processing summary");
    println!("├──────────────────────────────────────────┤");
    println!("│ {:<15} : {:>22} │", "Status", status);
    println!("│ {:<15} : {:>22} │", "Files moved", moved_count);
    println!("│ {:<15} : {:>22} │", "Folders created", folder_count);
    println!("│ {:<15} : {:>22} │", "Errors", error_count);
    println!("│ {:<15} : {:>22} │", "Duration", duration);
    println!("│ {:<15} : {:>22} │", "Throughput", throughput);
    println!("└──────────────────────────────────────────┘");

    Ok(())
}
