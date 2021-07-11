use structopt::StructOpt;
use toml::value::Array;
use toml::value::Table;
use serde::Deserialize;
use home;
use std::io::{Write};
use std::process::{Command, Stdio, Output};

/// Constants
const SANDMAN_DIR: &str = "Sandman";

/// Command arguments
#[derive(Debug, StructOpt)]
struct Args {
    action: String,
    container_name: String,
}

/// Build related configuration of a container
#[derive(Debug, Deserialize)]
struct ContainerConfigBuild {
    image_name: String,
    instructions: String,
}

/// Run related configuration of a container
#[derive(Debug, Deserialize)]
struct ContainerConfigRun {
    x11: bool,
    dri: bool,
    ipc: bool,
    pulseaudio: bool,
    dbus: bool,
    net: bool,
    volumes: Array,
    devices: Array,
}

/// The configuration of a container
#[derive(Debug, Deserialize)]
struct ContainerConfig {
    build: ContainerConfigBuild,
    run: ContainerConfigRun,
    env: Table,
}

/// A container is represented here
#[derive(Debug)]
struct Container {
    name: String,
    file: String,
    config: ContainerConfig,
}

/// Loads and returns a container based on the TOML configuration
fn get_container(container_name: &String) -> Container {
    let config_filename = format!("{}/{}/{}.toml", home::home_dir().unwrap().display(), SANDMAN_DIR, container_name);
    let config_raw = std::fs::read_to_string(&config_filename).unwrap();
    let config: ContainerConfig = toml::from_str(&config_raw).unwrap();

    Container {
        name: container_name.clone(),
        file: config_filename,
        config: config,
    }
}

/// Builds a given container
fn build_container(container: &Container) -> Result<Output, Output> {
    let image_name = container.name.clone();
    let dockerfile = container.config.build.instructions.clone();
    let build_arguments = vec!["bud", "-f", "-", "-t", &image_name];

    println!("Building {}", image_name);
    println!("Dockerfile Instructions:\n{}", dockerfile);

    dbg!(&image_name);
    dbg!(&dockerfile);
    dbg!(&build_arguments);

    // Set std file descriptors with pipes because we need to work with them
    let mut buildah = Command::new("buildah")
        .args(&build_arguments)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .unwrap();

    // Pass the dockerfile instructions via stdin
    let mut stdin = buildah.stdin.take().expect("Failed to open stdin");
    std::thread::spawn(move || {
        stdin.write_all(dockerfile.as_bytes()).expect("Failed to write to stdin")
    });

    // Wait command to finish and capture result
    let output = buildah.wait_with_output().expect("Failed to read stdout");
    let stdout = String::from_utf8_lossy(&output.stdout);
    let stderr = String::from_utf8_lossy(&output.stderr);
    dbg!(&stdout);
    dbg!(&stderr);

    if output.status.success() {
        return Ok(output);
    }
    else {
        return Err(output);
    }
}

/// Runs a given container
fn run_container(container: &Container) -> Result<(), Output> {
    println!("Running container...");
    // TODO
    dbg!(&container.name);
    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::from_args();
    let container = get_container(&args.container_name);

    dbg!(&args);
    dbg!(&container);

    if args.action == "run" {
        match run_container(&container) {
            Ok(output) => {
                println!("Container spawned successfully");
            },
            Err(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout);
                let stderr = String::from_utf8_lossy(&output.stderr);
                println!("Failed to run container: {} {}", stdout, stderr);
            },
        }
    }
    else if args.action == "build" {
        match build_container(&container) {
            Ok(output) => {
                println!("Image built successfully");
            },
            Err(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout);
                let stderr = String::from_utf8_lossy(&output.stderr);
                println!("Failed to build image: {} {}", stdout, stderr);
            },
        }
    }
    else {
        panic!("Action {} invalid", &args.action);
    };

    Ok(())
}
