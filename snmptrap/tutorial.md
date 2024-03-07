

# Building a Log Receiver

This assumes that you have already followed the [Building a receiver tutorial](https://opentelemetry.io/docs/collector/building/receiver/) previously.

Here's how OpenTelemetry thinks of a [log](https://opentelemetry.io/docs/concepts/signals/logs/)

    A log is a timestamped text record, either structured (recommended) or unstructured, with metadata.

In order to implement a logs receiver you will need the following:

* A `Config` implementation to enable the log receiver to gather and validate its configurations within the Collector’s `config.yaml`.
* A `receiver.Factory` implementation so the Collector can properly instantiate the log receiver component.
* A `LoggerReceiver` implementation that is responsible to collect the telemetry, convert it to the internal log representation, and hand the information to the next consumer in the pipeline.


In this tutorial we will create a sample log receiver called `snmptrapreceiver` that receives SNMP traps and generates logs from the trap. The next sections will guide you through the process of implementing the steps above in order to create the receiver, so let’s get started.

# About SNMP Traps
The Simple Network Management Protocol (SNMP) defines both the wire protocol (eg a Protocol Data Unit (PDU)) between devices as well as the
human-readable description (eg a Management Information Base (MIB)). Networking devices record information locally as Object Identifiers (OIDs)
and SNMP pollers can read these OIDs and parse the results (ie this is what the `snmpreceiver` module does).

A network device can also send out an event (eg interface up/down, temperature alarm, threshold exceeded) to a program that can understand it (eg network management software).
Technically, there are two types of events:

* `traps`: an event that does not require an acknowledgement that it has been received 
* `informs`: an event that does require an acknowledgement that it has been received 

SNMP is typically based on UDP, but can be configured to run on TCP.

For simplicity, anywhere you see the word `traps` assume that it could read `traps` or `informs`.

# Logs vs Traces
Whenever you create a receiver, you have to ask yourself "What is the most appropriate use for the collected data?"

For SNMP traps, the opinion I'm going to choose is to treat the SNMP trap as a structured log record. We should be able to recreate
our trap on the exporter side so that it can be possible to send out traps as traps again. For instance, we could receive an SNMP v3 `inform`
from a remote device and then send it out to a network monitoring tool as a SNMP v1 `trap` without loss of information.

Another use case is to be able to capture our traps and store them locally on the machine, so that we could replay these traps
for quality assurance or load-testing reasons. For quality assurance, these captured traps could be used in Collector unit tests as well
as to exercise code and processing on external systems.

In a sense, SNMP traps are also more like events than traces: a point in time happening that doesn't represent a transaction or activity
across multiple systems. And when we use these network events, we're more likely to want to process them more like log events than traces.

Also, there's already a tutorial for traces, but not for logs.


# Copying an Existing Receiver
We will use much of the configuration from the [snmpreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snmpreceiver) as
a basis for our own receiver. This will help to ensure consistency when trying to maintain configuration.

Specifically, we will re-use much of the `config.go` to maintain this consistency.

FIXME: import the final config.go

## About `mdatagen`
Note that in our code, we're going to need to update the `internal/metadata/generated_status.go` file to account for receiving logs vs traces.

As explained in the [documentation](https://github.com/open-telemetry/opentelemetry-collector/blob/main/cmd/mdatagen/README.md) we need to install
the OpenTelemetry metadata generator into our path in order to use it.

    # You've already done this part in the preceding tutorials
    % git clone https://github.com/open-telemetry/opentelemetry-collector.git
    % cd opentelemetry-collector
    % cd cmd/mdatagen && go install .

Now we need to update the `metadata.yaml` file with our name and capabilities.

The complete schema for [`metadata.yaml`](https://github.com/open-telemetry/opentelemetry-collector/blob/main/cmd/mdatagen/metadata-schema.yaml) provides
for information about optional and mandatory components. For our logs example, all we really want is a meta-data definition (ie `LogsStability`).

Our final `metadata.yaml` file contents:

    type: snmptrap

    status:
      class: receiver
      stability:
        alpha: [logs]
      distributions: []
      codeowners:
        active: [me]


To compile our `generated_status.go` file:

    % mdatagen metadata.yaml
    % cat internal/metadata/generated_status.go

There isn't any output from the command.

Here's what our generated file looks like:


    // Code generated by mdatagen. DO NOT EDIT.

    package metadata

    import (
	    "go.opentelemetry.io/collector/component"
	    "go.opentelemetry.io/otel/metric"
	    "go.opentelemetry.io/otel/trace"
    )

    var (
	    Type      = component.MustNewType("snmptrap")
	    scopeName = "go.opentelemetry.io/collector"
    )

    const (
	    LogsStability = component.StabilityLevelAlpha
    )

    func Meter(settings component.TelemetrySettings) metric.Meter {
	    return settings.MeterProvider.Meter(scopeName)
    }

    func Tracer(settings component.TelemetrySettings) trace.Tracer {
	    return settings.TracerProvider.Tracer(scopeName)
    }


