positionFoot = ->
    body = document.body
    html = document.documentElement
    wheight = Math.max body.offsetHeight,
        html.clientHeight, html.offsetHeight
    foot = document.getElementById "foot"
    height = foot.height || foot.offsetHeight
    foot.style.top = "#{wheight - height - 10}px"

addEvent = (elem, type, eventHandle) ->
    return if not elem?
    if elem.addEventListener
        elem.addEventListener type, eventHandle, false
    else if elem.attachEvent
        elem.attachEvent "on#{type}", eventHandle
    else
        elem["on#{type}"] = eventHandle

addEvent window, "load", positionFoot
addEvent window, "resize", positionFoot
