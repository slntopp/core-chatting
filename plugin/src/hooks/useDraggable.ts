interface MakeDraggableOprions {
    direction: 'H'
    resizer: HTMLElement
    first: HTMLElement
    second: HTMLElement
}

export default () => {
    function makeDraggable({resizer, first, second, direction = 'H'}: MakeDraggableOprions) {
        let coordinates: any;
        resizer.onmousedown = onMouseDown;

        function onMouseDown(e: any) {
            coordinates = {
                e,
                offsetLeft: resizer.offsetLeft,
                offsetTop: resizer.offsetTop,
                firstWidth: first!.offsetWidth,
                secondWidth: second!.offsetWidth
            };

            document.onmousemove = onMouseMove;
            document.onmouseup = () => {
                document.onmousemove = document.onmouseup = null;
            }
        }

        function onMouseMove({clientY, clientX}: MouseEvent) {
            const {e: eCoordinates, firstWidth, secondWidth, offsetLeft} = coordinates

            const delta = {
                x: clientX - eCoordinates.clientX,
                y: clientY - eCoordinates.clientY
            };

            if (direction === "H") // Horizontal
            {
                delta.x = Math.min(Math.max(delta.x, -firstWidth),
                    secondWidth);

                resizer.style.left = offsetLeft + delta.x + "px";
                first!.style.width = (firstWidth + delta.x) + "px";
            }
        }
    }

    return {makeDraggable}
}