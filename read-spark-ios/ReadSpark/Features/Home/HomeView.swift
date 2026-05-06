import SwiftUI

struct HomeView: View {
    @StateObject private var viewModel = HomeViewModel()
    @State private var selectedArticle: ArticleSummary?

    var body: some View {
        NavigationView {
            List {
                if viewModel.isLoading {
                    ProgressView()
                        .frame(maxWidth: .infinity, alignment: .center)
                }

                if let error = viewModel.error {
                    Text(error)
                        .foregroundColor(.red)
                }

                ForEach(viewModel.articles) { article in
                    ArticleCard(article: article) {
                        selectedArticle = article
                    }
                    .listRowSeparator(.hidden)
                }
            }
            .listStyle(.plain)
            .navigationTitle("ReadSpark")
            .sheet(item: $selectedArticle) { article in
                NavigationView {
                    ReaderView(articleId: article.id)
                }
            }
        }
    }
}
