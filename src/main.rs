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

struct ToggleImplication {
    env: Vec<String>,
    volumes: Vec<String>,
    devices: Vec<String>,
    args: Vec<String>,
}

struct Toggles {
    x11: ToggleImplication,
    dri: ToggleImplication,
    ipc: ToggleImplication,
    pulseaudio: ToggleImplication,
    dbus: ToggleImplication,
    net: ToggleImplication,
}

impl ContainerConfigRun {
    fn to_args(&self) -> Vec<String> {
        let toggles = get_toggles();
        let mut volumes: Vec<String> = vec![];
        let mut devices: Vec<String> = vec![];
        let mut env: Vec<String> = vec![];
        let mut args: Vec<String> = vec![];
        let mut arguments: Vec<String> = vec![];

        // Default arguments
        arguments.extend(vec![
            String::from("--interactive"),
            String::from("--tty"),
            String::from("--rm"),
        ]);

        println!("Converting toggles into arguments");

        if self.x11 {
            volumes.extend(toggles.x11.volumes);
            devices.extend(toggles.x11.devices);
            env.extend(toggles.x11.env);
            args.extend(toggles.x11.args);
        }
        if self.dri {
            volumes.extend(toggles.dri.volumes);
            devices.extend(toggles.dri.devices);
            env.extend(toggles.dri.env);
            args.extend(toggles.dri.args);
        }
        if self.ipc {
            volumes.extend(toggles.ipc.volumes);
            devices.extend(toggles.ipc.devices);
            env.extend(toggles.ipc.env);
            args.extend(toggles.ipc.args);
        }
        if self.pulseaudio {
            volumes.extend(toggles.pulseaudio.volumes);
            devices.extend(toggles.pulseaudio.devices);
            env.extend(toggles.pulseaudio.env);
            args.extend(toggles.pulseaudio.args);
        }
        if self.dbus {
            volumes.extend(toggles.dbus.volumes);
            devices.extend(toggles.dbus.devices);
            env.extend(toggles.dbus.env);
            args.extend(toggles.dbus.args);
        }
        if self.net {
            volumes.extend(toggles.net.volumes);
            devices.extend(toggles.net.devices);
            env.extend(toggles.net.env);
            args.extend(toggles.net.args);
        }

        //volumes.extend(self.volumes);
        //env.extend(self.env);
        //devices.extend(self.devices);

        for volume in volumes.iter() {
            arguments.extend(vec![String::from("--volume"), String::from(volume)]);
        }
        for device in devices.iter() {
            arguments.extend(vec![String::from("--device"), String::from(device)]);
        }
        for env_ in env.iter() {
            arguments.extend(vec![String::from("--env"), String::from(env_)]);
        }
        for arg in args.iter() {
            arguments.push(String::from(arg));
        }

        arguments
    }
}

fn get_toggles() -> Toggles {
    let x11 = ToggleImplication {
        env: vec![String::from(format!("DISPLAY={}", env!("DISPLAY")))],
        volumes: vec![String::from("/tmp/.X11-unix:/tmp/.X11-unix")],
        devices: vec![],
        args: vec![],
    };

    let dri = ToggleImplication {
        env: vec![],
        volumes: vec![],
        devices: vec![String::from("/dev/dri")],
        args: vec![],
    };

    let ipc = ToggleImplication {
        env: vec![],
        volumes: vec![],
        devices: vec![],
        args: vec![String::from("--ipc"), String::from("host")],
    };

    let pulseaudio = ToggleImplication {
        env: vec![],
        volumes: vec![
            String::from("/etc/machine-id:/etc/machine-id:ro"),
            String::from(format!("{}/pulse/native:{}/pulse/native", env!("XDG_RUNTIME_DIR"), env!("XDG_RUNTIME_DIR"))),
        ],
        devices: vec![],
        args: vec![],
    };

    let dbus = ToggleImplication {
        env: vec![String::from(format!("DBUS_SESSION_BUS_ADDRESS=unix:path={}/bus", env!("XDG_RUNTIME_DIR")))],
        volumes: vec![String::from(format!("{}/bus:{}/bus", env!("XDG_RUNTIME_DIR"), env!("XDG_RUNTIME_DIR")))],
        devices: vec![],
        args: vec![],
    };

    let net = ToggleImplication {
        env: vec![],
        volumes: vec![],
        devices: vec![],
        args: vec![String::from("--network"), String::from("slirp4netns")],
    };

    Toggles {
        x11: x11,
        dri: dri,
        ipc: ipc,
        pulseaudio: pulseaudio,
        dbus: dbus,
        net: net,
    }
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
    let args = container.config.run.to_args();
    dbg!(args);
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
