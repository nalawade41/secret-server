#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DeployApi } from '../lib/deploy-api';

const { CDK_DEFAULT_ACCOUNT, CDK_DEFAULT_REGION } = process.env;

const app = new cdk.App();
new DeployApi(app, 'DeployStack', {
    env: { account: CDK_DEFAULT_ACCOUNT, region: CDK_DEFAULT_REGION },
});