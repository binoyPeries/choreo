# React Starter Web Application - Container Image

This sample demonstrates how to deploy a web application as a container image in Choreo without the source code.

## Deploying in Choreo
The following command will create the component, deployment track, deployable artifact and the deployment in Choreo.

```bash
choreoctl apply -f samples/deploying-applications/use-prebuilt-image/react-spa-webapp/react-starter.yaml
```

## Retrieving the Invocation URL

You can retrieve the invocation URL of the deployment by running the following command.

```bash
choreoctl get endpoint --organization default-org --project default-project --component  react-starter-image
```

This will output the endpoints in the `default-org` namespace. For this specific sample
you will see an endpoint with the name starting with `react-starter-image-deployment-webapp-<hash>`. 
It should have the address as `https://react-starter-image-development.choreo.localhost`

## Accessing the Web Application

In order to access the web application, you can use one of the approaches mentioned in [README.md](../../../../README.md#7-test-the-deployed-webapp).

For this sample, we will use kubectl port-forward to access the web application.

1. Run the following command to port-forward the gateway.

    ```bash
    kubectl port-forward svc/envoy-choreo-system-gateway-external-<hash> -n choreo-system 4430:443
    ```
   
2. Access the web application from your browser using the following URL.
    - https://react-starter-image-development.choreo.localhost:4430
