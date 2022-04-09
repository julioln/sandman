use home::{home_dir};
use crate::SANDMAN_STORAGE_DIR;
use crate::args::Args;
use crate::args::ExecuteArgs;
use crate::toggles::Toggles;

use serde::Deserialize;
use std::io::{Write};
use std::process::{Command, Stdio, ExitStatus};
use std::fs::create_dir_all;
use std::path::Path;

/// Build related configuration of a container
#[derive(Debug, Deserialize)]
pub struct ContainerConfigBuild {
    instructions: String,
}

/// Run related configuration of a container, expected in the TOML file
#[derive(Debug, Deserialize)]
pub struct ContainerConfigRun {
    #[serde(default)]
    x11: bool,

    #[serde(default)]
    wayland: bool,

    #[serde(default)]
    dri: bool,

    #[serde(default)]
    ipc: bool,

    #[serde(default)]
    pulseaudio: bool,

    #[serde(default)]
    dbus: bool,

    // Kept for backwards compatibility
    #[serde(default)]
    net: bool,

    #[serde(default)]
    network: String,

    #[serde(default)]
    uidmap: bool,

    #[serde(default)]
    volumes: Vec<String>,

    #[serde(default)]
    home: bool,

    #[serde(default)]
    devices: Vec<String>,

    #[serde(default)]
    env: Vec<String>,

    #[serde(default)]
    ports: Vec<String>,

    #[serde(default)]
    name: String,

    #[serde(default)]
    memory_limit: String,

    #[serde(default)]
    args: Vec<String>,
}

/// The configuration of a container
#[derive(Debug, Deserialize)]
pub struct ContainerConfig {
    pub build: ContainerConfigBuild,
    pub run: ContainerConfigRun,
}

/// A container is represented here
#[derive(Debug)]
pub struct Container {
    pub name: String,
    pub basename: String,
    pub file: String,
    pub config: ContainerConfig,
}

/// The main container object
impl Container {

    /// Returns a vector of podman arguments compiled from all configuration
    pub fn running_args(&self) -> Vec<String> {
        let cli_args = Args::cli_args();
        let toggles = Toggles::get_toggles();
        let mut volumes: Vec<String> = vec![];
        let mut devices: Vec<String> = vec![];
        let mut ports: Vec<String> = vec![];
        let mut env: Vec<String> = vec![];
        let mut args: Vec<String> = vec![];
        let mut arguments: Vec<String> = vec![];

        // Default arguments
        arguments.extend(vec![
            String::from("run"),
            String::from("--hostname"),
            self.name.clone().replace('/', "_"),
            String::from("--interactive"),
            String::from("--tty"),
        ]);

        if !cli_args.keep {
            arguments.extend(vec![
                String::from("--rm"),
            ]);
        }

        if !self.config.run.name.is_empty() {
            arguments.extend(vec![String::from("--name"), self.config.run.name.clone()])
        }

        if !self.config.run.memory_limit.is_empty() {
            arguments.extend(vec![String::from("--memory"), self.config.run.memory_limit.clone()])
        }

        if self.config.run.home {
            let local_home_mount = format!("{}/{}/{}",
                home_dir().unwrap().display(),
                SANDMAN_STORAGE_DIR,
                self.basename
            );
            create_dir_all(Path::new(&local_home_mount)).unwrap();
            volumes.extend([
                format!("{}:/home", local_home_mount)
            ]);
        }

        // Collect all configuration from toggles that are enabled
        // TODO Improve this maybe with a hash
        if self.config.run.x11 {
            volumes.extend(toggles.x11.volumes);
            devices.extend(toggles.x11.devices);
            env.extend(toggles.x11.env);
            args.extend(toggles.x11.args);
        }
        if self.config.run.wayland {
            volumes.extend(toggles.wayland.volumes);
            devices.extend(toggles.wayland.devices);
            env.extend(toggles.wayland.env);
            args.extend(toggles.wayland.args);
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
        if self.config.run.uidmap {
            volumes.extend(toggles.uidmap.volumes);
            devices.extend(toggles.uidmap.devices);
            env.extend(toggles.uidmap.env);
            args.extend(toggles.uidmap.args);
        }

        if self.config.run.net {
            volumes.extend(toggles.net.volumes);
            devices.extend(toggles.net.devices);
            env.extend(toggles.net.env);
            args.extend(toggles.net.args);
        }
        else if self.config.run.network.is_empty() {
            arguments.extend(vec![String::from("--network"), String::from("none")]);
        }
        else {
            arguments.extend(vec![String::from("--network"), self.config.run.network.clone()]);
        }

        // Add customized configuration
        volumes.extend(self.config.run.volumes.clone());
        env.extend(self.config.run.env.clone());
        devices.extend(self.config.run.devices.clone());
        ports.extend(self.config.run.ports.clone());
        args.extend(self.config.run.args.clone());

        // Mix command line overrides
        if !cli_args.env.is_empty() {
            env.extend(cli_args.env.clone());
            env.sort();
            env.dedup();
        }

        for volume in volumes.iter() {
            arguments.extend(vec![String::from("--volume"), String::from(volume)]);
        }
        for device in devices.iter() {
            arguments.extend(vec![String::from("--device"), String::from(device)]);
        }
        for port in ports.iter() {
            arguments.extend(vec![String::from("-p"), String::from(port)]);
        }
        for env_ in env.iter() {
            arguments.extend(vec![String::from("--env"), String::from(env_)]);
        }
        for arg in args.iter() {
            arguments.push(String::from(arg));
        }

        // Image name
        arguments.push(self.name.clone());

        // Pass extra arguments after the image name, if we have any coming from
        // the command line
        if let Some(other) = cli_args.execute {
            match other {
                ExecuteArgs::Other(extra_args) => arguments.extend(extra_args),
            }
        }

        arguments
    }

    /// Builds a given container
    pub fn build(&self) -> Result<ExitStatus, ExitStatus> {
        let cli_args = Args::cli_args();
        let image_name = self.name.clone();
        let dockerfile = self.config.build.instructions.clone();
        let mut build_arguments: Vec<String> = vec![String::from("bud")];

        if cli_args.cache {
            build_arguments.extend(vec![
                String::from("--layers=true"),
            ]);
        }

        build_arguments.extend(vec![
            String::from("-f"),
            String::from("-"),
            String::from("-t"),
            String::from(&image_name),
        ]);

        if cli_args.verbose {
            dbg!(&build_arguments);
        }

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
            Ok(status)
        }
        else {
            Err(status)
        }
    }

    /// Runs a given container
    pub fn run(&self) -> Result<ExitStatus, ExitStatus> {
        let args = self.running_args();
        let cli_args = Args::cli_args();

        if cli_args.verbose {
            dbg!(&args);
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
            Ok(status)
        }
        else {
            Err(status)
        }
    }

}
