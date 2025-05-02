# lilypad-concat-files-module

This module demonstrates file inputs to a Lilypad job. It concatenates text files with an optional separator.

### Run

Try this module using the `lilypad` CLI. See https://docs.lilypad.tech/lilypad/quickstart/cli for install and wallet setup instructions.

Create an `inputs` directory with `one.txt` and `two.txt` files with some text in them. Run the module to concatenate them:

```sh
lilypad run github.com/Lilypad-Tech/lilypad-concat-files-module:main --inputs-path=./inputs
```

You can also add `three.txt` and `four.txt` files, but they are optional.

Concatenate the files with a separator:

```sh
lilypad run github.com/Lilypad-Tech/lilypad-concat-files-module:main --inputs-path=./inputs -i Separator="=========="
```

### Using input files

Lilypad modules can declare required and optional file inputs to a job. The required inputs must be present or the job will be rejected. No other files beyond the required and optional files are permitted. Declared files must be unique across required and optional file lists.

Input files are declared in the Lilypad module template. Our module declares the following input files:

https://github.com/Lilypad-Tech/lilypad-concat-files-module/blob/875689fe956e18eb6223e5ed5b8bfd829e09bc55/lilypad_module.json.tmpl#L7-L10

Input files are available to a job at `/inputs` in the container where the job runs. Module Dockerfiles should create this directory to make sure it available.

https://github.com/Lilypad-Tech/lilypad-concat-files-module/blob/875689fe956e18eb6223e5ed5b8bfd829e09bc55/Dockerfile#L11

File inputs can be used with templated textual inputs to a job. These inputs are the common case used when passing prompts to an LLM. In our example, we define a `Separator` that is passed as an environment variable to the job.

https://github.com/Lilypad-Tech/lilypad-concat-files-module/blob/875689fe956e18eb6223e5ed5b8bfd829e09bc55/lilypad_module.json.tmpl#L24-L26

Once inside the job, use textual inputs from their environment variables and input files from the `/inputs` directory.

### Build

Build the docker image:

```sh
docker build -t lilypad-concat-files .
```

Publish the docker image:

```sh
docker buildx build --platform linux/amd64,linux/arm64 -t bgins/lilypad-concat-files-module:latest --push .
```
