/**
 * Created by dashie on 08/05/2017.
 */

$(document).on("click", "#add_new_form", function() {
    console.log("adding a new file fields");
    // Clone the file template
    let new_file = $("div.files > div").last().clone();

    // Empty the inputs
    new_file.find("input[type='text']").val("");
    new_file.find("textarea[type='text']").val("");

    // Update the index
    let index = parseInt(new_file.data("id"));
    let newIndex = index + 1;

    // This only works in debugger, not here, thank you JS
    //new_file.data("id", newIndex);
    new_file.attr("data-id", newIndex);

    new_file.attr("class", "form_" + newIndex + " gitxt_new_file");
    new_file.find("label[for='file_"+index+"_filename']").attr("for", "file_"+newIndex+"_filename");
    new_file.find("label[for='file_"+index+"_content']").attr("for", "file_"+newIndex+"_content");

    new_file.find("input[id='file_"+index+"_filename']").attr("id", "file_"+newIndex+"_filename");
    new_file.find("textarea[id='file_"+index+"_content']").attr("id", "file_"+newIndex+"_content");

    new_file.find("div[class='file_idx']").html("file " + newIndex);

    // Append to the files
    new_file.appendTo("div.files");
});

$(document).on("click", ".btn-delete-file", function() {
    console.log("removing file fields");
    let parent = $(this).closest(".gitxt_new_file");
    if ($("div.files > div").length > 1) {
        parent.remove();
    } else {
        $(this).tooltip({title: "not for this one", placement: "right", container: 'body'});
        $(this).tooltip('show');
        console.log("nope I won't do that");
    }
});