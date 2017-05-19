hljs.configure({useBR: false});

$('pre#code').each(function(i, block) {
    hljs.highlightBlock(block);
});

$(document).on("click", "a.delete_link", function(e) {
    e.preventDefault();

    let $this = $(this);

    let rOwner = $this.data("owner");
    let rHash = $this.data("hash");

    if (confirm("Are you sure ?") === true) {
        $.post($this.data('url'), {
            "_csrf": csrf,
            "owner": rOwner,
            "hash": rHash
            }).done(function (data) {
                //window.location.href = data.redirect;
            });
    }
});

$(document).on("click", "button.img-loader", function(e) {
    e.preventDefault();

    let imgUrl = $(this).data('src');
    $(this).parent().html("<img src='" + imgUrl + "' alt='gitxt image' />");
});
