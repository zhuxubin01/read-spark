import SwiftUI

struct HomeView: View {
    @StateObject private var viewModel = HomeViewModel()

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
                        // Task 7 will wire reader navigation.
                    }
                    .listRowSeparator(.hidden)
                }
            }
            .listStyle(.plain)
            .navigationTitle("ReadSpark")
        }
    }
}
