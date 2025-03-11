# Xcode test shard calculation

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-xcode-test-shard-calculation?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-xcode-test-shard-calculation/releases)

Calculates the test shards for your Xcode tests.

<details>
<summary>Description</summary>

This step looks at the test bundle and calculates the test shards for your Xcode tests. It finds all of the tests and
divides them into the specified number of shards.

These shards can be used to run the tests in parallel, which can significantly speed up the testing process. Use the
[Xcode Test without building](https://www.bitrise.io/integrations/steps/xcode-test-without-building) step's `Test Selection`
input to specify which test shard information to use.

### Related Steps
- [Xcode Test without building](https://www.bitrise.io/integrations/steps/xcode-test-without-building)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `product_path` | The path of the test bundle.  The step supports the following formats: - xcresrun - xctestproducts  It will use the specified file to collect the built tests and generate the test shards. | required |  |
| `shard_count` | The number of test shards to create.  The output folder will contain `shard_count` number of files, each containing the tests to run in that shard. | required |  |
| `shard_calculation` | Defines the strategy to use when splitting the tests into shards  The available options are: - `alphabetically`: The tests are sorted alphabetically and split into shards | required | `alphabetically` |
| `destination` | Destination specifier describes the device to use as a destination.  The input value sets xcodebuild's `-destination` option.  In a CI environment, a Simulator device called `Bitrise iOS default` is already created. It is a compatible device with the selected Simulator runtime, pre-warmed for better performance.  If a device with this name is not found (e.g. in a local dev environment), the first matching device will be selected. | required | `platform=iOS Simulator,name=Bitrise iOS default,OS=latest` |
| `verbose` | Enable logging additional information for debugging. | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_TEST_SHARDS_PATH` | This folder contains the generated test shard information. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-xcode-test-shard-calculation/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-xcode-test-shard-calculation/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
