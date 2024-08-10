import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {LogGroup, RetentionDays} from "aws-cdk-lib/aws-logs";
import {
  Cors,
  LambdaIntegration,
  LogGroupLogDestination,
  MethodLoggingLevel,
  RestApi,
  RestApiProps
} from "aws-cdk-lib/aws-apigateway";
import {
  Runtime,
  Tracing,
  Function,
  Code,
  RuntimeManagementMode, LayerVersion,
} from "aws-cdk-lib/aws-lambda";
import * as path from "node:path";
import {Stack} from "aws-cdk-lib";
import {Effect, Policy, PolicyStatement} from "aws-cdk-lib/aws-iam";


export class DeployApi extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const ENV = 'prod';
    const apiGatewayName = `${ENV}-api-gateway`;
    console.log(apiGatewayName);
    // Create a CloudWatch Log Group
    const logGroup = new LogGroup(this, `${ENV}-sso-logs`, {
      retention: RetentionDays.INFINITE,
      logGroupName: `/aws/apigateway/${ENV}-api-logs`,
    });


    // Create the Lambda function
    const apiHandler = new Function(this, "api-sls", {
      runtime: Runtime.PROVIDED_AL2,
      code: Code.fromAsset(path.join(__dirname, "..", "..", "bin")),
      environment: {
        ENV: ENV!,
        READ_TIMEOUT: "5s",
        WRITE_TIMEOUT: "5s",
        MAX_HEADER_BYTES: "1048576",
        DB_TABLE_NAME: "secrets",
      },
      tracing: Tracing.ACTIVE,
      memorySize: 512,
      functionName: `${ENV}-api`,
      timeout: cdk.Duration.seconds(30),
      handler: "bootstrap", // TODO name does not matter when the runtime is PROVIDED_AL2
      runtimeManagementMode: RuntimeManagementMode.AUTO,
    });

    // Set up the Lambda function to use the AWS X-Ray tracing
    const arch = 'amd64';
    const otelCollectorVersion = '0-102-1:1';
    const otelAWSAccount = '901920570463';
    const otelCollectorArn = `arn:aws:lambda:${Stack.of(this).region}:${otelAWSAccount}:layer:aws-otel-collector-${arch}-ver-${otelCollectorVersion}`
    const layerName = `${ENV}-api-otel-layer`;
    apiHandler.addLayers(
        LayerVersion.fromLayerVersionArn(this, layerName, otelCollectorArn)
    );

    const policy = new Policy(this, "secret-policy", {
      statements: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: [
            "appconfig:*",
            'dynamodb:*',
          ],
          resources: ["*", 'arn:aws:dynamodb:*:*:table/*'],
        }),
      ],
    });

    apiHandler.role?.attachInlinePolicy(policy);

    // Configure API Gateway properties
    const apiGatewayProps: RestApiProps = {
      description: `${ENV} API Gateway`,
      deploy: true,
      deployOptions: {
        tracingEnabled: true,
        stageName: ENV,
        loggingLevel: MethodLoggingLevel.INFO,
        accessLogDestination: new LogGroupLogDestination(logGroup),
      },
      cloudWatchRole: true,
      restApiName: apiGatewayName,
      defaultCorsPreflightOptions: {
        allowMethods: Cors.ALL_METHODS,
        allowHeaders: [
          "Content-Type",
          "Authorization",
          "X-Amz-Date",
          "X-Api-Key",
          "X-Amz-Security-Token",
          "X-Amz-User-Agent",
          "X-Public-Id",
          "Accept",
        ],
        allowOrigins: Cors.ALL_ORIGINS,
        maxAge: cdk.Duration.days(1),
      },
      defaultIntegration: new LambdaIntegration(apiHandler, {
        proxy: true,
      }),
    };

    // Create a new instance of RestApi
    const apiGateway = new RestApi(this, apiGatewayName, apiGatewayProps);

    apiGateway.root.addProxy({
      defaultIntegration: new LambdaIntegration(apiHandler, {
        proxy: true,
      }),
      anyMethod: true, // "true" is the default
    });

    new cdk.CfnOutput(this, 'ApiGatewayUrl', {
      value: apiGateway.url,
      description: 'The base URL of the API Gateway',
      exportName: `${ENV}-api-gateway-url`,
    });
  }
}
