use structopt::StructOpt;

/// Command arguments
#[derive(Debug, PartialEq, StructOpt)]
pub struct Args {
    #[structopt(short, long)]
    pub verbose: bool,

    #[structopt(short, long)]
    pub keep: bool,

    #[structopt(short, long)]
    pub cache: bool,

    pub action: String,
    pub container_name: String,

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
