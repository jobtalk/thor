config.Version = 1.0;

var family = "family name";
var targetGroupArn = "何らかのARN";
var serviceName = "サービス名";
var taskDefinitionName = "task definition name";

service.body = {
    "Cluster": "backend",
    "DeploymentConfiguration": {
        "MaximumPercent": 200,
        "MinimumHealthyPercent": 100
    },
    "DesiredCount": 1,
    "LoadBalancers": [
        {
            "ContainerName": "nginx",
            "ContainerPort": 80,
            "TargetGroupArn": targetGroupArn
        }
    ],
    "Role": "ecsServiceRole",
    "ServiceName": serviceName,
    "TaskDefinition": taskDefinitionName
};

taskDefinition.body = {
    "ContainerDefinitions": [
        {
            "Cpu": 0,
            "Essential": true,
            "ExtraHosts": null,
            "Image": "ieee0824/nginx-template",
            "MemoryReservation": 128,
            "Name": "nginx",
            "Links": [
                "api"
            ],
            "PortMappings": [
                {
                    "HostPort": 0,
                    "ContainerPort": 80,
                    "Protocol": "tcp"
                }
            ]
        },
        {
            "Cpu": 0,
            "Essential": true,
            "ExtraHosts": null,
            "Image": "ieee0824/dummy-app",
            "MemoryReservation": 512,
            "Name": "api"
        }
    ],
    "Family": family,
    "NetworkMode": "bridge"
}


