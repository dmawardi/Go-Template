// Grab the checkbox action form and add an event listener: commitMultiAction
// document.getElementById('deleteForm').addEventListener('submit', commitMultiAction);

// Select all the rows in the table
function selectAllRows(checkbox) {
  // Grab all the checkboxes
  var rowCheckboxes = document.querySelectorAll(".select-row");
  //   Iterate through the checkboxes and set checked to the value of the select-all checkbox
  for (var i = 0; i < rowCheckboxes.length; i++) {
    rowCheckboxes[i].checked = checkbox.checked;
  }
}

// Update the select-all checkbox status whenever a row checkbox is clicked
function updateSelectAll(checkbox) {
  // Grab the select-all checkbox
  var selectAllCheckbox = document.getElementById("select-all");
  //   If the row checkbox is not checked, then the select-all checkbox should not be checked
  if (!checkbox.checked) {
    selectAllCheckbox.checked = false;
  } else {
    //  If the row checkbox is checked, then check if all the row checkboxes are checked
    var allRowCheckboxesChecked = true;
    // Grab all the checkboxes
    var rowCheckboxes = document.querySelectorAll(".select-row");
    // Iterate through the checkboxes and set allRowCheckboxesChecked to false if any of the row checkboxes are not checked
    for (var i = 0; i < rowCheckboxes.length; i++) {
      if (!rowCheckboxes[i].checked) {
        allRowCheckboxesChecked = false;
        break;
      }
    }
    // Set the select-all checkbox to the value of allRowCheckboxesChecked
    selectAllCheckbox.checked = allRowCheckboxesChecked;
  }
}

// Sort the table by the given column
function sortTable(orderBy) {
  // Init value
  let orderValue = "";
  // Grab the url parameters
  const urlParams = new URLSearchParams(window.location.search);
  // Check current order parameter
  const currentOrderBy = urlParams.get("order");
  // If the current order parameter is the same as the given column, then reverse the order
  if (currentOrderBy === orderBy) {
    // If the current order parameter is ascending, then set the order parameter to descending
    orderValue = orderBy + "_desc";
    urlParams.set("order", orderValue);
  } else {
    orderValue = orderBy;
    // Otherwise, set the order parameter to the given column
    urlParams.set("order", orderValue);
  }
  // Reload the page with the new url parameters
  window.location.search = urlParams;
}

async function commitMultiAction(event, schemaDeleteUrl) {
  // Prevent default behavior of submission
  event.preventDefault();

  // Collect selected user IDs
  const selectedItems = [];
  // Select all the checkboxes that are checked and iterate through them
  document
    .querySelectorAll('input[name="selected_items"]:checked')
    .forEach(function (item) {
      // Push the value of the checkbox to the selectedItems array
      selectedItems.push(item.value);
    });
  console.log("selectedItems: ", selectedItems);
  console.log("schemaDeleteUrl: ", schemaDeleteUrl);

  try {
    selectedItemsJson = JSON.stringify({ selected_items: selectedItems });
    console.log("selectedItemsJson: ", selectedItemsJson);
    // Send a DELETE request to the server
    const response = await fetch(schemaDeleteUrl + "/bulk-delete", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      // Convert the selectedItems array to JSON and send it in the body of the request
      body: selectedItemsJson,
    });

    if (!response.ok) {
      throw new Error("Something went wrong");
    }
    console.log("success");
    const data = await response.json();
    console.log(data);
    // Handle success
  } catch (error) {
    console.error("Error:", error);
    // Handle errors
  }
}
