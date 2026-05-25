# Agritech Edge Sensor Daemon & Action Service

A lightweight edge computing agent built with Go, acting as the nervous system for smart farm operations. It simulates environmental sensors, ensures local hardware autonomy, and interfaces seamlessly with centralized cloud monitoring.

## Key Features

* **Edge Autonomy (SPOF Prevention):** Executes critical control logic (e.g., sprinklers, ventilation) locally, ensuring greenhouse safety even if the network to the control plane is disconnected.
* **Telemetry Ingestion:** Continuously streams mock environmental metrics (Temperature, Humidity, Soil Moisture) to a centralized InfluxDB.
* **Action Webhook Receiver:** Listens for Grafana alerts from the Monitor-Node to trigger physical actuators (e.g., Virtual Water Pump).
* **Secure Credentials:** Utilizes environment variables to prevent hardcoded tokens in the source code.

## Edge Architecture Flow

```mermaid
graph TD
    subgraph Farm-Node [Farm-Node / Edge Environment]
        SD[Sensor Daemon] -->|Sensor Readings| LCL[Local Control Logic]
        LCL -->|Instant Trigger| Actuators[Sprinkler / Ventilation]
        
        AS[Action Service Webhook] -->|Grafana Alert| Pump[Water Pump Action]
    end

    subgraph Monitor-Node [Monitor-Node / Control Plane]
        SD -->|Push Telemetry| IDB[(InfluxDB)]
        Grafana[Grafana Alerting] -->|HTTP POST| AS
    end
    
    style Farm-Node fill:#e8f4f8,stroke:#333,stroke-width:2px
    style Monitor-Node fill:#f9eaeb,stroke:#333,stroke-width:2px
