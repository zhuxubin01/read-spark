import SwiftUI
import UIKit

struct ReaderTextView: UIViewRepresentable {
    let text: String
    let onLongPressWord: (String) -> Void

    func makeUIView(context: Context) -> UITextView {
        let textView = UITextView()
        textView.isEditable = false
        textView.isSelectable = true
        textView.font = UIFont.systemFont(ofSize: 18)
        textView.textColor = UIColor.label
        textView.backgroundColor = UIColor.systemBackground

        let gesture = UILongPressGestureRecognizer(target: context.coordinator, action: #selector(Coordinator.handleLongPress(_:)))
        gesture.minimumPressDuration = 0.5
        textView.addGestureRecognizer(gesture)

        return textView
    }

    func updateUIView(_ uiView: UITextView, context: Context) {
        uiView.text = text
    }

    func makeCoordinator() -> Coordinator {
        Coordinator(self)
    }

    class Coordinator: NSObject {
        let parent: ReaderTextView

        init(_ parent: ReaderTextView) {
            self.parent = parent
        }

        @objc func handleLongPress(_ gesture: UILongPressGestureRecognizer) {
            guard gesture.state == .began else { return }

            guard let textView = gesture.view as? UITextView else { return }
            let point = gesture.location(in: textView)

            if let textPosition = textView.closestPosition(to: point),
               let wordRange = textView.tokenizer.rangeEnclosingPosition(textPosition, with: .word, inDirection: UITextDirection.storage(.forward)) {
                let word = textView.text(in: wordRange) ?? ""
                if !word.isEmpty {
                    parent.onLongPressWord(word)
                }
            }
        }
    }
}
