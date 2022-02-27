use structopt::StructOpt;

/// Command arguments
#[derive(Debug, PartialEq, StructOpt)]
pub struct Args {
    /// Prints debugging information
    #[structopt(short, long)]
    pub verbose: bool,

    /// RUN: Keep the container (do not include --rm)
    #[structopt(short, long)]
    pub keep: bool,

    /// RUN: Override or set environment variables
    #[structopt(short, long)]
    pub env: Vec<String>,

    /// BUILD: Use layer cache (--layers=true)
    #[structopt(short, long)]
    pub cache: bool,

    /// The action: run, build
    pub action: String,

    /// The container name to be found in SANDMAN_DIR
    pub container_name: String,

    /// Optional arguments to pass in the run action
    #[structopt(subcommand)]
    pub execute: Option<ExecuteArgs>,
}

/// Optional subcommand arguments when running a container
#[derive(Debug, PartialEq, StructOpt)]
pub enum ExecuteArgs {
    #[structopt(external_subcommand)]
    Other(Vec<String>),
}

impl Args {
    /// Command line arguments
    pub fn cli_args() -> Args {
        Args::from_args()
    }
}
