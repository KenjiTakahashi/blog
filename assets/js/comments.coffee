width = $("#content").width()
container = $("#comments")

animateIn = (data) ->
    $("#content").animate width: "#{width - 220}px", 1000
    $("#title").animate width: "#{width - 220}px", 1000
    $("canvas").animate width: "#{width - 230}px", 1000
    container.animate width: "462px", 1000, ->
        if data?
            $data = $(data)
            $data.hide()
            container.append $data
        height = "#{$(window).height() - 20}px"
        $("#dsq1").css height: height
        $("#disqus_thread").show()
        $(@).animate height: height, 1000, ->
            $("a", @).attr 'id', 'acommentsa'

$("#comments a").toggle ->
    if container.children().length < 2
        path = $(@).attr('href').split '/'
        path = path[path.length - 1]
        $.ajax "/comments/#{path}", {
            beforeSend: ->
                $.spin 'commentsa', @
            success: animateIn
            complete: $.unspin
        }
    else
        animateIn()
, ->
    $(@).removeAttr 'id'
    container.animate height: "40px", 1000, ->
        $("#disqus_thread").hide()
        $(@).animate width: "100px", 1000, ->
            $("a", @).attr 'id', 'commentsa'
        $("#content").animate width: "#{width}px", 1000
        $("#title").animate width: "#{width}px", 1000
        $("canvas").animate width: "#{width - 10}px", 1000
