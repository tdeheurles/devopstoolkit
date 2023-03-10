import json
import time
import logging
import os
import process

from sys import exit
from typing import Dict, List
from termcolor import cprint

from cli.Component import Component, Command
from cli.log import cli_print
from cli.parameters import key, key_aws_region, \
    get_and_assert_additional_parameters, key_stack_name, key_prompt, key_disable_rollback
from cli.theme import theme_error0, theme_warn0


def deploy_cloudformation_stack(params: Dict[str, str], component: Component, command: Command) -> int:

    logger = logging.getLogger()

    # USAGE
    all_params = get_and_assert_additional_parameters(
        initial_params=params,
        usage=lambda cn=component.name, cd=command: component.print_command_usage(cn, cd),
        new_required_params=[
            {key: key_aws_region},
            {key: key_stack_name},
            {key: key_disable_rollback}
        ])

    assert_dependencies(["aws","sam"])

    # EXECUTION
    cloudformation_directory = f"{component.directory}/cloudformation"
    initial_template = f"{cloudformation_directory}/template.yaml"
    sam_build_directory = f"{cloudformation_directory}/.aws-sam/build"
    template_file = f"{sam_build_directory}/template.yaml"

    status_to_wait = [
        "CREATE_IN_PROGRESS",
        "CREATE_FAILED",
        "DELETE_FAILED",
        "DELETE_IN_PROGRESS",
        "REVIEW_IN_PROGRESS",
        "ROLLBACK_FAILED",
        "ROLLBACK_IN_PROGRESS",
        "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS",
        "UPDATE_IN_PROGRESS",
        "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS",
        "UPDATE_ROLLBACK_FAILED",
        "UPDATE_ROLLBACK_IN_PROGRESS",
        "IMPORT_IN_PROGRESS",
        "IMPORT_ROLLBACK_IN_PROGRESS",
        "IMPORT_ROLLBACK_FAILED"
    ]
    os.environ['AWS_DEFAULT_REGION'] = all_params[key_aws_region]
    while True:
        aws_cf_list_process, aws_cf_list_stdout, aws_cf_list_stderr= process.run(
            [
                f"aws cloudformation list-stacks"
                f"  --region={all_params[key_aws_region]}"
            ])
        if aws_cf_list_process.returncode != 0:
            logger.error("Unable to run aws cloudformation list-stacks")
            exit(1)

        aws_cf_list = json.loads(aws_cf_list_stdout)
        stack = [
            stack
            for stack in aws_cf_list["StackSummaries"]
            if stack["StackName"] == all_params[key_stack_name]
            if stack["StackStatus"] != "DELETE_COMPLETE"
        ]
        logger.info(stack)
        if len(stack) == 0:
            logger.info(f'no stack {all_params[key_stack_name]} exists. Creating ...')
            break
        if stack[0]["StackStatus"] in status_to_wait:
            logger.info(f'Stack is in status {stack["StackStatus"]}. Waiting ...')
            time.sleep(60)
        else:
            break

        # 2020 - April
        # CREATE_COMPLETE                              Successful creation of one or more stacks.
        # CREATE_IN_PROGRESS                           Ongoing creation of one or more stacks.
        # CREATE_FAILED                                Unsuccessful creation of one or more stacks. View the stack events to see any associated error messages. Possible reasons for a failed creation include insufficient permissions to work with all resources in the stack, parameter values rejected by an AWS service, or a timeout during resource creation.
        # DELETE_COMPLETE                              Successful deletion of one or more stacks. Deleted stacks are retained and viewable for 90 days.
        # DELETE_FAILED                                Unsuccessful deletion of one or more stacks. Because the delete failed, you might have some resources that are still running; however, you can't work with or update the stack. Delete the stack again or view the stack events to see any associated error messages.
        # DELETE_IN_PROGRESS                           Ongoing removal of one or more stacks.
        # REVIEW_IN_PROGRESS                           Ongoing creation of one or more stacks with an expected StackId but without any templates or resources.
        # ROLLBACK_COMPLETE                            Successful removal of one or more stacks after a failed stack creation or after an explicitly canceled stack creation. Any resources that were created during the create stack operation are deleted.
        #                                              This status exists only after a failed stack creation. It signifies that all operations from the partially created stack have been appropriately cleaned up. When in this state, only a delete operation can be performed.
        # ROLLBACK_FAILED                              Unsuccessful removal of one or more stacks after a failed stack creation or after an explicitly canceled stack creation. Delete the stack or view the stack events to see any associated error messages.
        # ROLLBACK_IN_PROGRESS                         Ongoing removal of one or more stacks after a failed stack creation or after an explicitly canceled stack creation.
        # UPDATE_COMPLETE                              Successful update of one or more stacks.
        # UPDATE_COMPLETE_CLEANUP_IN_PROGRESS          Ongoing removal of old resources for one or more stacks after a successful stack update. For stack updates that require resources to be replaced, CloudFormation creates the new resources first and then deletes the old resources to help reduce any interruptions with your stack. In this state, the stack has been updated and is usable, but CloudFormation is still deleting the old resources.
        # UPDATE_FAILED                                Unsuccessful update of one or more stacks. View the stack events to see any associated error messages.
        # UPDATE_IN_PROGRESS                           Ongoing update of one or more stacks.
        # UPDATE_ROLLBACK_COMPLETE                     Successful return of one or more stacks to a previous working state after a failed stack update.
        # UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS Ongoing removal of new resources for one or more stacks after a failed stack update. In this state, the stack has been rolled back to its previous working state and is usable, but CloudFormation is still deleting any new resources it created during the stack update.
        # UPDATE_ROLLBACK_FAILED                       Unsuccessful return of one or more stacks to a previous working state after a failed stack update. When in this state, you can delete the stack or continue rollback. You might need to fix errors before your stack can return to a working state. Or, you can contact AWS Support to restore the stack to a usable state.
        # UPDATE_ROLLBACK_IN_PROGRESS                  Ongoing return of one or more stacks to the previous working state after failed stack update.
        # IMPORT_IN_PROGRESS                           The import operation is currently in progress.
        # IMPORT_COMPLETE                              The import operation successfully completed for all resources in the stack that support resource import.
        # IMPORT_ROLLBACK_IN_PROGRESS                  Import will roll back to the previous template configuration.
        # IMPORT_ROLLBACK_FAILED                       The import rollback operation failed for at least one resource in the stack. Results will be available for the resources CloudFormation successfully imported.
        # IMPORT_ROLLBACK_COMPLETE                     Import successfully rolled back to the previous template configuration.

    # BUILD
    sam_build, _, _ = process.run(
        [
            f"sam build"
            f"  --template {initial_template}"
            f"  --build-dir {sam_build_directory}"
        ], logger)

    if sam_build.returncode != 0:
        logger.error("Unable to run sam build command")
        exit(1)

    # DEPLOY
    os.environ['AWS_DEFAULT_REGION'] = all_params[key_aws_region]
    rollback = '--disable-rollback' if all_params[key_disable_rollback] == 'true' else '--no-disable-rollback'
    sam_deploy, _, _ = process.run(
        [
            f"sam deploy"
            f"  --template-file {template_file}"
            f"  --stack-name {all_params[key_stack_name]}"
            f"  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM"
            f"  --no-fail-on-empty-changeset"
            f"  {rollback}"
        ], logger)

    if sam_deploy.returncode != 0:
        logger.error("Unable to run sam build command")
        exit(1)

    return 0


def assert_dependencies(dependencies: List[str]):
    for command in dependencies:
        _, stdout, stderr = process.run([f'command -v {command}'])
        if stdout == "":
            cprint(f"Error, command {command} is required", theme_error0)
            exit(1)
