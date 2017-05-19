hljs.configure({useBR: false});

$('pre#code').each(function(i, block) {
    hljs.highlightBlock(block);
});

// Chain all on-click actions on document
$(document).on("click", "a.delete_link", function(e) {
    // Delete link actions : Ajax with CSRF
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
                // TODO FIXME
                //window.location.href = data.redirect;
            });
    }
}).on("click", "button.img-loader", function(e) {
    // img-loader will load the image inplace
    e.preventDefault();

    let imgUrl = $(this).data('src');
    $(this).parent().html("<img src='" + imgUrl + "' alt='gitxt image' />");
}).on("click", "button.pdf-loader", function(e) {
    // pdf-loader will load PDF.js viewer with the PDF inplace
   e.preventDefault();

   let pdfUrl = $(this).data('src');
   $(this).parent().html('<iframe width="100%" height="600px" src="/pdfjs-1.7.225/web/viewer.html?file=' + pdfUrl + '"></iframe>');
});