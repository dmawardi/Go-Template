// URL to receive requests on server
const server_port = "8080";
const server_address = "http://localhost:" + server_port;
const policyEditUrl = server_address + "/admin/policy/";

// Function to edit policy from form submission
async function editPolicy(e, actionToComplete, role, resource, action) {
  console.log(
    "Parameters:\nAction: " +
      actionToComplete +
      "\nRole: " +
      role +
      "\nResource: " +
      resource +
      "\nAction: " +
      action
  );
  //   Prevent default behavior
  e.preventDefault();

  //   Depending on action to complete, send request to server
  switch (actionToComplete) {
    case "create":
      // Send Post request to add policy
      response = await sendRequest("post", {
        role: role,
        resource: resource,
        action: action,
      });
      console.log("response in add: ", response);
      break;

    case "delete":
      response = await sendRequest("delete", {
        role: role,
        resource: resource,
        action: action,
      });
      console.log("response in delete: ", response);
      break;
    default:
      response = false;
  }

  console.log("response: ", response);
  if (!response) {
    alert("Action failed.");
  } else {
    // If successful
    // Once complete, Reload the page
    window.location.reload();
  }
}

// Takes requested action in form: "POST", "DELETE" and sends requested action to server
// Takes data in form of JSON {role: role, resource: resource, action: action} and sends requested action to server
async function sendRequest(requestedAction, policyData) {
  try {
    // Convert requested action (POST/DELETE) to all caps
    capitalizedAction = requestedAction.toUpperCase();
    //   Build a param by replacing "/" with "-"
    slugParam = policyData.resource.replace(/\//g, "-");

    //  Build request URL
    requestUrl = policyEditUrl + slugParam;

    console.log("capitalizedAction: ", capitalizedAction);
    console.log("slugParam: ", slugParam);
    console.log("requestUrl: ", requestUrl);

    // Convert the selectedItems array to JSON
    policyDataJson = JSON.stringify(policyData);
    console.log("policyDataJson: ", policyDataJson);
    // Send a DELETE request to the server
    const response = await fetch(requestUrl, {
      method: capitalizedAction,
      headers: {
        "Content-Type": "application/json",
      },
      // Convert the selectedItems array to JSON and send it in the body of the request
      body: policyDataJson,
    });

    // If the response is not ok, then throw an error
    if (!response.ok) {
      return false;
    }
    console.log("success. Response: ", response);
    // Else convert json response to data
    const data = await response.json();
    console.log("data: ", data);
    return data;
  } catch (error) {
    console.error("Error:", error);
    return false;
  }
}
