

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
The Simple Network Management Protocol (SNMP) defines both the wire protocol (eg a Protocol Data Unit) between devices as well as the
human-readable description (eg a Management Information Base (MIB)). Networking devices record information locally as Object Identifiers (OIDs)
and SNMP pollers can read these OIDs and parse the results (ie this is what the `snmpreceiver` module does).

A network device can also send out an event (eg interface up/down, temperature alarm, threshold exceeded) to a program that can understand it (eg network management software).
Technically, there are two types of events:

* `traps`: an event that does not require an acknowledgement that it has been received 
* `informs`: an event that does require an acknowledgement that it has been received 

SNMP is typically based on UDP, but can be configured to run on TCP.

# Copying an Existing Receiver
We will use much of the configuration from the [snmpreceiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/snmpreceiver) as
a basis for our own receiver. This will help to ensure consistency when trying to maintain configuration.

Specifically, we will re-use much of the `config.go` to maintain this consistency.

FIXME: import the final config.go

