width = $("#content").width()
container = $("#comments")

animateIn = ->
    $("#content").animate width: "#{width - 220}px", 1500
    $("#title").animate width: "#{width - 220}px", 1500
    $("canvas").animate width: "#{width - 230}px", 1500
    container.animate width: "462px", 1500, ->
        $("#disqus_thread").show()
        $(this).animate height: $(window).height(), 1500

$("#comments a").toggle ->
    if container.children().length < 2
        jQuery.get "/comments", (data) ->
            container.append data
            $("#disqus_thread").hide()
            animateIn()
    else
        animateIn()
, ->
    container.animate height: "40px", 1500, ->
        $("#disqus_thread").hide()
        $(this).animate width: "100px", 1500
        $("#content").animate width: "#{width}px", 1500
        $("#title").animate width: "#{width}px", 1500
        $("canvas").animate width: "#{width - 10}px", 1500
