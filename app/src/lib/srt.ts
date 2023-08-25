interface SrtLine {
    id: string
    time: string
    text: string
}

export class Srt {
    lines: SrtLine[]

    constructor(srtString: string) {
        this.lines = []

        for(const lineStr of srtString.split("\n\n")) {
            const lineParts = lineStr.split("\n")

            const id = lineParts[0]
            const time = lineParts[1]
            const text = lineParts.slice(2).join()

            this.lines.push({ id, time, text })
        }
    }

    toString(): string {
        let srtString = ""
        for(const line of this.lines) {
            srtString += line.id + "\n"
            srtString += line.time + "\n"
            srtString += line.text
            srtString += "\n\n"
        }
        return srtString
    }
}
