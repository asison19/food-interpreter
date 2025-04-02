from google.cloud import artifactregistry
import asyncio


async def list_docker_images():
    """
        List docker images in GCP Artifact Registry
    """
    client = artifactregistry.ArtifactRegistryAsyncClient()

    # Initialize request argument(s)
    request = artifactregistry.ListDockerImagesRequest(
        #parent="us-central1-docker.pkg.dev/food-interpreter/food-interpreter-repository",
        parent="projects/food-interpreter/locations/us-central1/repositories/food-interpreter-repository",
    )

    # Make the request
    page_result = await client.list_docker_images(request=request)

    # Handle the response
    amount = 0
    async for response in page_result:
        print(response)
        amount += 1

    return page_result, amount

async def delete_old_docker_images(page_result, amount, max_amount=13):
    """
        Delete old docker images
    """
    client = artifactregistry.ArtifactRegistryAsyncClient()

    # order by upload time?


    for i in range(amount - max_amount):
        print(i)
        
        # find oldest image
        # TODO this is funky
        async for response in page_result:
            oldest_image = response
            break
        #oldest_image = next(iter(page_result))
        async for response in page_result:
            if response.upload_time < oldest_image.upload_time:
                oldest_image = response
                print('foobar')
        
        print("Deleting:", oldest_image)

        # Initialize request argument(s)
        request = artifactregistry.DeleteVersionRequest(
            name=oldest_image.uri,
        )

        operation = client.delete_version(request=request)
        print("Waiting for operation to complete...")
        response = (await operation).result()

    # Handle the response
    #print(response)

if __name__ == '__main__':
    #await list_docker_images()
    page_result, amount = asyncio.run(list_docker_images())
    asyncio.run(delete_old_docker_images(page_result, amount))
    print('done')
    

    #request = artifactregistry.DeletePackageRequest(
    #    name="name_value",
    #)

    ## Make the request
    #operation = client.batch_delete_versions(request=request)

    #print("Waiting for operation to complete...")

    #response = (await operation).result()

    ## Handle the response
    #print(response)