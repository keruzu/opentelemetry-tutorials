

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

# Pre-requisite Software Installation
In order to perform the testing later on, we'll need to ensure that we have the necessary SNMP command-line utilities installed.

On RHEL-type systems:

   % sudo dnf install net-snmp net-snmp-utils net-snmp-libs net-snmp-agent-libs



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

## Working with Go Workspaces
If you followed the previous tutorial, you will be using Go [Workspaces](https://go.dev/doc/tutorial/workspaces) so we'll need to do a few quick updates
before we can start testing or doing other work.

    % pwd
    /your/path/here/snmptrap
    % cd ..
    % go work use snmptrap
    % cd snmptrap
    % go mod init
    % go mod tidy

This will rebuild the `go.work` file and allow you to build and test the modules

# Building out a Sane `config.go` with Unit Testing
Here's our plan of action from here:

1. Rename all of the module packaes to `snmptrap`
1. Remove references to metrics, scraper helpers and anything that refers to polling
1. Rename the `Endpoint` configuration to instead talk about `ListenAddres` as we're not polling for data, we're receiving packets
1. Move the unit tests out of the way until `go build` works
1. Move the `config.go` unit test back in and fix things until `go test` works

I'm going to spare you the drudgery of the above and get us to this point.

Here's the contents of our `config.go`

FIXME: add the config.go contents here

# Adding `Start()` and `Shutdown()` Functions
Now we're ready to start listening!

Our `factory.go` contents showing the details:

FIXME: add the factory.go



## A Little More About SNMP Traps
There are a number of [PDU structures](https://cdpstudio.com/manual/cdp/snmpio/about-snmp.html) for SNMP communications.
[RFC 3416](https://www.ietf.org/rfc/rfc3416.txt) introduces the SNMP v2 trap PDU, for reference.

If you need an authoritative reference, this tutorial is not it.

Here's an outline of the key metadata that we have:

* PDU type
* Length
* Enterprise OID
* Agent Address (ie the thing sending the trap)
* Generic Trap Type
* Specific Trap number
* Time stamp
* List of varbinds (eg OIDs + values)

There are two trap types:

* Generic traps: `coldStart`, `warmStart`, `linkDown`, `linkUp`, `authenticationFailure`, and `egpNeighborLoss`
* Enterprise-specific traps: traps defined by the vendor of the device

A trap will provide with OID indices, which will look something like: `1.3.6.1.4.1. ....`
This represents a hierarchical tree of dots and numbers, with meaning according to the RFCs and assigned by a numbering authority.
A MIB can be used to convert the numbers to names, with the idea being that you can map numbers to names like DNS maps IP addresses to FQDN names.

Note that in theory everything is well-ordered, well-defined and everything can be uniquely determined in a "simple" way. For instance, each vendor
will have their own MIB that defines the data so that it can be uniquely compiled, and there are no name clashes, badly ordered names and everything is consistent.
Much like web browsers and HTML, there's a lot in the history of SNMP that can and has gone wrong in practice.

For our purposes, we'd like to capture the essential metadata from our traps, and create a key/value store with our OIDs indices and their values.
We can then pass this through another layer (eg an OTel Collector processor or an actual SNMP manager) to convert to human-friendly names and then
even later on be able to respond appropriately to the trap.


# About SNMP Trap Processing Performance
Note that the kernel UDP buffer ring stores UDP packets (eg our SNMP traps), and this is buffer is overwritten with new incoming UDP
packets as time goes on. Our goal is to be able to process the contents of this buffer as fast as possible so that we can actually
keep pace with the oncoming stream of SNMP packets.

There are a few things we can do:

* tune the kernel to request more UDP buffer space to increase the total size of the ring buffer
FIXME: indicate the Linux kernel parameter
* ensure that we can specify a big UDP buffer size so that we can hint to the kernel that we care about those packets
FIXME: verify the above statement
* ensure that our code is as performant as possible
* drop any packets we don't care about as soon as possible
* add a load balancer in front of this and spin up more containers, processes or whatever

There are options, but we are going to ingore them all at this point.

