use structopt::StructOpt;
use serde::Deserialize;
use home;
use std::io::{Write};
use std::process::{Command, Stdio, ExitStatus};

/// Constants
const SANDMAN_DIR: &str = "Sandman";

/// Command arguments
#[derive(Debug, StructOpt)]
struct Args {
    #[structopt(short, long)]
    verbose: bool,
    action: String,
    container_name: String,
    execute: String,
}

/// Build related configuration of a container
#[derive(Debug, Deserialize)]
struct ContainerConfigBuild {
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
    uidmap: bool,
    volumes: Vec<String>,
    devices: Vec<String>,
    env: Vec<String>,
}

/// The configuration of a container
#[derive(Debug, Deserialize)]
struct ContainerConfig {
    build: ContainerConfigBuild,
    run: ContainerConfigRun,
}

/// A container is represented here
#[derive(Debug)]
struct Container {
    name: String,
    file: String,
    config: ContainerConfig,
}

#[derive(Hash, Eq, PartialEq, Debug)]
struct ToggleImplication {
    env: Vec<String>,
    volumes: Vec<String>,
    devices: Vec<String>,
    args: Vec<String>,
}

#[derive(Hash, Eq, PartialEq, Debug)]
struct Toggles {
    x11: ToggleImplication,
    dri: ToggleImplication,
    ipc: ToggleImplication,
    pulseaudio: ToggleImplication,
    dbus: ToggleImplication,
    net: ToggleImplication,
    uidmap: ToggleImplication,
}

impl Container {
    fn running_args(&self) -> Vec<String> {
        let toggles = get_toggles();
        let mut volumes: Vec<String> = vec![];
        let mut devices: Vec<String> = vec![];
        let mut env: Vec<String> = vec![];
        let mut args: Vec<String> = vec![];
        let mut arguments: Vec<String> = vec![];

        // Default arguments
        arguments.extend(vec![
            String::from("run"),
            String::from("--hostname"),
            String::from(self.name.clone().replace("/", "_")),
            String::from("--interactive"),
            String::from("--tty"),
            String::from("--rm"),
        ]);

        if self.config.run.x11 {
            volumes.extend(toggles.x11.volumes);
            devices.extend(toggles.x11.devices);
            env.extend(toggles.x11.env);
            args.extend(toggles.x11.args);
        }
        if self.config.run.dri {
            volumes.extend(toggles.dri.volumes);
            devices.extend(toggles.dri.devices);
            env.extend(toggles.dri.env);
            args.extend(toggles.dri.args);
        }
        if self.config.run.ipc {
            volumes.extend(toggles.ipc.volumes);
            devices.extend(toggles.ipc.devices);
            env.extend(toggles.ipc.env);
            args.extend(toggles.ipc.args);
        }
        if self.config.run.pulseaudio {
            volumes.extend(toggles.pulseaudio.volumes);
            devices.extend(toggles.pulseaudio.devices);
            env.extend(toggles.pulseaudio.env);
            args.extend(toggles.pulseaudio.args);
        }
        if self.config.run.dbus {
            volumes.extend(toggles.dbus.volumes);
            devices.extend(toggles.dbus.devices);
            env.extend(toggles.dbus.env);
            args.extend(toggles.dbus.args);
        }
        if self.config.run.net {
            volumes.extend(toggles.net.volumes);
            devices.extend(toggles.net.devices);
            env.extend(toggles.net.env);
            args.extend(toggles.net.args);
        }
        if self.config.run.uidmap {
            volumes.extend(toggles.uidmap.volumes);
            devices.extend(toggles.uidmap.devices);
            env.extend(toggles.uidmap.env);
            args.extend(toggles.uidmap.args);
        }

        volumes.extend(self.config.run.volumes.clone());
        env.extend(self.config.run.env.clone());
        devices.extend(self.config.run.devices.clone());

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

        arguments.push(self.name.clone());

        arguments
    }

    /// Builds a given container
    fn build(&self) -> Result<ExitStatus, ExitStatus> {
        let image_name = self.name.clone();
        let dockerfile = self.config.build.instructions.clone();
        let build_arguments = vec!["bud", "-f", "-", "-t", &image_name];

        // Set stdin with pipe because we need to pass the dockerfile using it
        let mut buildah = Command::new("buildah")
            .args(&build_arguments)
            .stdin(Stdio::piped())
            .stdout(Stdio::inherit())
            .stderr(Stdio::inherit())
            .spawn()
            .unwrap();

        // Pass the dockerfile instructions via stdin
        let mut stdin = buildah.stdin.take().expect("Failed to open stdin");
        std::thread::spawn(move || {
            stdin.write_all(dockerfile.as_bytes()).expect("Failed to write to stdin")
        });

        let status = buildah.wait().expect("Failed to read stdout");

        if status.success() {
            return Ok(status);
        }
        else {
            return Err(status);
        }
    }

    /// Runs a given container
    fn run(&self) -> Result<ExitStatus, ExitStatus> {
        let mut args = self.running_args();
        let cli_args = cli_args();

        if cli_args.verbose {
            dbg!(&args);
        }

        if !cli_args.execute.is_empty() {
            args.push(cli_args.execute);
        }

        let mut podman = Command::new("podman")
            .args(&args)
            .stdin(Stdio::inherit())
            .stdout(Stdio::inherit())
            .stderr(Stdio::inherit())
            .spawn()
            .unwrap();

        let status = podman.wait().expect("Failed to read stdout");

        if status.success() {
            return Ok(status);
        }
        else {
            return Err(status);
        }
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
        env: vec![String::from(format!("XDG_RUNTIME_DIR={}", env!("XDG_RUNTIME_DIR")))],
        volumes: vec![
            String::from("/etc/machine-id:/etc/machine-id:ro"),
            String::from(format!("{}/pulse/native:{}/pulse/native", env!("XDG_RUNTIME_DIR"), env!("XDG_RUNTIME_DIR"))),
        ],
        devices: vec![],
        args: vec![],
    };

    let uidmap = ToggleImplication {
        env: vec![],
        volumes: vec![],
        devices: vec![],
        args: vec![
            String::from("--uidmap"),
            String::from("1000:0:1"),
            String::from("--uidmap"),
            String::from("0:1:1000"),
            String::from("--uidmap"),
            String::from("1001:1001:64536"),
            String::from("--user"),
            String::from("1000"),
        ],
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
        uidmap: uidmap,
    }
}

/// Loads and returns a container based on the TOML configuration
fn get_container(container_name: &String) -> Container {
    let config_filename = format!("{}/{}/{}.toml", home::home_dir().unwrap().display(), SANDMAN_DIR, container_name);
    let config_raw = std::fs::read_to_string(&config_filename).unwrap();
    let config: ContainerConfig = toml::from_str(&config_raw).unwrap();

    Container {
        name: format!("sandman/{}", container_name.clone()),
        file: config_filename,
        config: config,
    }
}

fn cli_args() -> Args {
   Args::from_args()
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = cli_args();
    let container = get_container(&args.container_name);

    if args.verbose {
        dbg!(&args);
        dbg!(&container);
    }

    if args.action == "run" {
        match container.run() {
            Err(status) => {
                println!("Failed to run container. Exit status: {}", status);
            },
            _ => {},
        }
    }
    else if args.action == "build" {
        match container.build() {
            Err(status) => {
                println!("Failed to build container. Exit status: {}", status);
            },
            _ => {},
        }
    }
    else if args.action == "args" {
        let arguments = container.running_args();
        let joined = arguments.join(" ");
        println!("{}", joined);
    }
    else {
        panic!("Action {} invalid", &args.action);
    };

    Ok(())
}
