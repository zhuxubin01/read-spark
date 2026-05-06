import SwiftUI

struct ReaderView: View {
    let articleId: UUID
    @StateObject private var viewModel = ReaderViewModel()
    @Environment(\.dismiss) private var dismiss
    @State private var selectedWord: String?

    var body: some View {
        VStack(spacing: 0) {
            if viewModel.isLoading {
                ProgressView()
                    .padding()
            } else if let error = viewModel.error {
                Text(error)
                    .foregroundColor(.red)
                    .padding()
            } else if let article = viewModel.article {
                ScrollView {
                    VStack(alignment: .leading, spacing: 16) {
                        Text(article.title)
                            .font(.title2)
                            .fontWeight(.bold)

                        Text("\(article.category) · \(article.difficulty) · \(article.wordCount) words")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        ReaderTextView(text: article.content) { word in
                            selectedWord = word
                        }
                        .frame(minHeight: 300)

                        if viewModel.showTranslation, let translation = article.translation {
                            Divider()
                            Text("译文")
                                .font(.headline)
                                .foregroundColor(.accentColor)
                            Text(translation)
                                .font(.body)
                                .foregroundColor(.secondary)
                        }
                    }
                    .padding()
                }
            }

            HStack {
                Button(action: {}) {
                    Image(systemName: "speaker.wave.2")
                }
                Spacer()
                Button(action: { viewModel.showTranslation.toggle() }) {
                    Image(systemName: "translate")
                }
                Spacer()
                Button(action: {}) {
                    Image(systemName: "square.and.pencil")
                }
                Spacer()
                Button(action: {}) {
                    Image(systemName: "bookmark")
                }
            }
            .padding()
            .background(Color(.systemBackground))
            .shadow(radius: 2)
        }
        .navigationTitle("阅读")
        .navigationBarTitleDisplayMode(.inline)
        .onAppear {
            viewModel.loadArticle(id: articleId)
        }
        .sheet(item: $selectedWord) { word in
            WordPopupView(word: word)
        }
    }
}

extension String: Identifiable {
    public var id: String { self }
}

struct WordPopupView: View {
    let word: String
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationView {
            VStack {
                Text(word)
                    .font(.largeTitle)
                    .padding()

                Text("[点击查词功能需要接入词典API]")
                    .foregroundColor(.secondary)

                Spacer()

                Button("加入生词本") {
                    dismiss()
                }
                .buttonStyle(.borderedProminent)
                .padding()
            }
            .navigationBarItems(trailing: Button("关闭") {
                dismiss()
            })
        }
    }
}
