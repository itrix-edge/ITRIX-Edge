Overview of Edge Appliance
==========================

[TOC]

**Note:** This document explain high level project overview, technical detail and implementation topics please refer `Specification.md`.

# Our Vision

Edge solution for appliance: with integrated software & reference hardware design, built for following goal:

1. Easy Installation
2. Easy Configuration
3. Easy Deployment
4. Easy Management

In the short term, we try to create an reasonable price micro GPU cluster, for edge computing applications. 
In this scope, we define an ethernet connected micro clusters inside the appliance, with extendable GPU enabled compute node from 3 to 24+; node management via Kubernetes; use custom software for easy application deployment; provide remote access design to make sure lower IT services in daily operation.

Wth the long term goal, we try to make a inter-appliance and federation enabled edge, for application level load dispatch. In our edge cases, all distributed appliances in different place will be a one virtual appliance in application level, with seamless worker load dispatch & failover feature enabled.

# Applications

This edge solution focus on following applications:

* Single appliance configuration:
	1. In-field AI inference application
	2. General micro services application
	3. HA required but cost sensitive roles

* Multiple appliance integration configuration:
	1. Multi-location service deployment
	2. Physical handover required service & application

